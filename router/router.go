package router

import (
	"sespima_api/internal/handler"
	"sespima_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth       *handler.AuthHandler
	User       *handler.UserHandler
	Kegiatan   *handler.KegiatanHandler
	Assessment *handler.AssessmentHandler
	Reward     *handler.RewardHandler
	Punishment *handler.PunishmentHandler
	Assignment *handler.AssignmentHandler
	WS         *handler.WSHandler
	Serdik     *handler.SerdikHandler
	Dashboard  *handler.DashboardHandler
	RP         *handler.RewardPunishmentHandler
}

func Register(r *gin.Engine, h *Handlers, jwtSecret string) {
	auth := middleware.Auth(jwtSecret)

	// ── Public auth routes ──────────────────────────────────────────────
	pub := r.Group("/api")
	pub.POST("/login", h.Auth.Login)
	pub.POST("/reset_password", h.Auth.ResetPassword)
	pub.POST("/verify_reset_token", h.Auth.VerifyResetToken)
	pub.POST("/refresh_token", h.Auth.RefreshToken)

	// ── Authenticated routes ────────────────────────────────────────────
	api := r.Group("/api", auth)

	// Auth
	api.GET("/me", h.Auth.GetMe)
	api.POST("/update_password", h.Auth.UpdatePassword)
	api.POST("/logout", h.Auth.Logout)
	api.POST("/generate_reset_token",
		middleware.RequireRoles("operator"),
		h.Auth.GenerateResetToken,
	)

	// Profile photo (authenticated user's own photo)
	api.POST("/me/photo", h.User.UploadProfilePhoto)
	api.DELETE("/me/photo", h.User.DeleteProfilePhoto)

	// Serdik self-profile (role-guarded)
	serdik := api.Group("/serdik", middleware.RequireRoles("serdik"))
	serdik.GET("/me", h.Serdik.GetMe)
	serdik.PUT("/me", h.Serdik.UpdateMe)

	// Static file serving for uploaded profile photos
	r.Static("/uploads", "./uploads")

	// Users / location
	users := api.Group("/users")
	users.GET("", h.User.GetAll)
	users.GET("/locations", h.User.GetAllLocations)
	users.PUT("/me/location", auth, h.User.UpdateMyLocation)
	users.GET("/:id/location", h.User.GetLocation)
	users.PUT("/:id/location", h.User.UpdateLocation)

	// WebSocket — auth is via ?token= query param, not header
	r.GET("/api/ws/locations", h.WS.StreamLocations)

	// Kegiatan (korsis creates, all authenticated can read)
	kg := api.Group("/kegiatan")
	kg.GET("", h.Kegiatan.GetAll)
	kg.GET("/rutin", h.Kegiatan.GetRutin)
	kg.GET("/today", h.Kegiatan.GetToday)
	kg.GET("/today/zones", h.Kegiatan.GetTodayZones)
	kg.POST("", middleware.RequireRoles("korsis", "operator"), h.Kegiatan.Create)
	kg.PUT("/:id", middleware.RequireRoles("korsis", "operator"), h.Kegiatan.Update)
	kg.DELETE("/:id", middleware.RequireRoles("korsis", "operator"), h.Kegiatan.Delete)

	// Attendance
	ab := api.Group("/absensi")
	ab.POST("", h.Kegiatan.CheckIn)
	ab.GET("/me", auth, h.Kegiatan.GetMyAttendance)
	ab.GET("/serdik/:serdikId", h.Kegiatan.GetAttendanceBySerdik)
	ab.GET("/kegiatan/:kegiatanId", h.Kegiatan.GetAttendanceByKegiatan)

	// Izin (leave requests)
	izin := api.Group("/izin")
	izin.POST("", middleware.RequireRoles("serdik"), h.Kegiatan.SubmitIzin)
	izin.GET("/serdik/:serdikId", h.Kegiatan.GetIzinBySerdik)
	izin.GET("/pending",
		middleware.RequireRoles("korsis", "patun", "operator"),
		h.Kegiatan.GetPendingIzin,
	)
	izin.PUT("/:id/review",
		middleware.RequireRoles("korsis", "patun", "operator"),
		h.Kegiatan.ReviewIzin,
	)

	// Akademik assessment
	ak := api.Group("/akademik")
	ak.GET("", h.Assessment.GetAkademik)
	ak.GET("/serdik/:serdikId", h.Assessment.GetAkademikBySerdik)
	ak.POST("", middleware.RequireRoles("gadik", "operator"), h.Assessment.CreateAkademik)
	ak.PUT("/:id", middleware.RequireRoles("gadik", "operator"), h.Assessment.UpdateAkademik)
	ak.DELETE("/:id", middleware.RequireRoles("gadik", "operator"), h.Assessment.DeleteAkademik)

	// Jasmani assessment (previously unrouted — now fixed)
	jas := api.Group("/jasmani")
	jas.GET("", h.Assessment.GetJasmani)
	jas.GET("/serdik/:serdikId", h.Assessment.GetJasmaniWithNilai)
	jas.POST("", middleware.RequireRoles("gadik", "operator"), h.Assessment.CreateJasmani)
	jas.PUT("/:id", middleware.RequireRoles("gadik", "operator"), h.Assessment.UpdateJasmani)
	jas.DELETE("/:id", middleware.RequireRoles("gadik", "operator"), h.Assessment.DeleteJasmani)

	// Mental assessment (previously unrouted — now fixed)
	men := api.Group("/mental")
	men.GET("", h.Assessment.GetMental)
	men.GET("/serdik/:serdikId", h.Assessment.GetMentalWithNilai)
	men.POST("", middleware.RequireRoles("korsis", "operator"), h.Assessment.CreateMental)
	men.PUT("/:id", middleware.RequireRoles("korsis", "operator"), h.Assessment.UpdateMental)
	men.DELETE("/:id", middleware.RequireRoles("korsis", "operator"), h.Assessment.DeleteMental)

	// Sociometry
	soc := api.Group("/sociometry")
	soc.GET("/periods", h.Assessment.GetPeriods)
	soc.POST("/periods",
		middleware.RequireRoles("korsis", "operator"),
		h.Assessment.CreatePeriod,
	)
	soc.PUT("/periods/:id/close",
		middleware.RequireRoles("korsis", "operator"),
		h.Assessment.ClosePeriod,
	)
	soc.POST("/evaluations", middleware.RequireRoles("serdik"), h.Assessment.SubmitEvaluation)

	// Health monitoring
	health := api.Group("/health")
	health.GET("/serdik/:serdikId", h.Assessment.GetHealthBySerdik)
	health.PUT("/serdik/:serdikId",
		middleware.RequireRoles("medis", "operator"),
		h.Assessment.UpsertHealthData,
	)
	health.POST("/records",
		middleware.RequireRoles("medis"),
		h.Assessment.AddHealthRecord,
	)

	// Rewards — Maker-Checker: patun/gadik create (forced pending),
	// only korsis/operator approve or reject.
	rw := api.Group("/user_rewards")
	rw.GET("", h.Reward.GetAll)
	rw.GET("/pending", middleware.RequireRoles("korsis", "operator"), h.Reward.GetPending)
	rw.POST("", middleware.RequireRoles("patun", "gadik", "operator"), h.Reward.Create)
	rw.PUT("/:id", middleware.RequireRoles("patun", "gadik", "operator"), h.Reward.Update)
	rw.DELETE("/:id", middleware.RequireRoles("patun", "gadik", "operator"), h.Reward.Delete)
	rw.PUT("/:id/approve", middleware.RequireRoles("korsis", "operator"), h.Reward.Approve)
	rw.PUT("/:id/reject", middleware.RequireRoles("korsis", "operator"), h.Reward.Reject)

	// Punishments — Maker-Checker: patun/gadik create (forced pending),
	// only korsis/operator approve or reject.
	pun := api.Group("/punishment_logs")
	pun.GET("", h.Punishment.GetAll)
	pun.GET("/pending", middleware.RequireRoles("korsis", "operator"), h.Punishment.GetPending)
	pun.GET("/user/:user_id", h.Punishment.GetByUserID)
	pun.POST("", middleware.RequireRoles("patun", "gadik", "operator"), h.Punishment.Create)
	pun.PUT("/:id", middleware.RequireRoles("patun", "gadik", "operator"), h.Punishment.Update)
	pun.DELETE("/:id", middleware.RequireRoles("patun", "gadik", "operator"), h.Punishment.Delete)
	pun.PUT("/:id/approve", middleware.RequireRoles("korsis", "operator"), h.Punishment.Approve)
	pun.PUT("/:id/reject", middleware.RequireRoles("korsis", "operator"), h.Punishment.Reject)

	// Aspect-based (Mental Kepribadian) reward/punishment Maker-Checker.
	// Catalog is readable by makers + checker; records are created by makers.
	makerCheckerRoles := []string{"patun", "gadik", "korsis", "operator"}
	rp := api.Group("/rp")
	rp.GET("/rules", middleware.RequireRoles(makerCheckerRoles...), h.RP.GetRules)
	rp.GET("/serdik/search", middleware.RequireRoles(makerCheckerRoles...), h.RP.SearchSerdik)
	rp.POST("/records", middleware.RequireRoles("patun", "gadik", "operator"), h.RP.CreateRecord)
	rp.GET("/records/pending", middleware.RequireRoles("korsis", "operator"), h.RP.GetPending)
	rp.GET("/records/serdik/:serdikId", middleware.RequireRoles(makerCheckerRoles...), h.RP.GetBySerdik)
	rp.GET("/records/me", middleware.RequireRoles("serdik"), h.RP.GetMyRecords)

	// Serdik dashboard (aggregated scores, rewards, punishments, history)
	api.GET("/dashboard/serdik", middleware.RequireRoles("serdik"), h.Dashboard.GetSerdikDashboard)

	// Unified Korsis approval inbox (pending rewards + punishments + izin)
	api.GET("/inbox/korsis", middleware.RequireRoles("korsis", "operator"), h.Dashboard.GetKorsisInbox)

	// Assignments (tugas)
	asgn := api.Group("/assignments")
	asgn.GET("", h.Assignment.GetAll)
	asgn.GET("/:id", h.Assignment.GetByID)
	asgn.POST("", middleware.RequireRoles("gadik", "operator"), h.Assignment.Create)
	asgn.PUT("/:id", middleware.RequireRoles("gadik", "operator"), h.Assignment.Update)
	asgn.DELETE("/:id", middleware.RequireRoles("gadik", "operator"), h.Assignment.Delete)
	asgn.GET("/:id/submissions", h.Assignment.GetSubmissions)
	asgn.POST("/:id/submissions", middleware.RequireRoles("serdik"), h.Assignment.Submit)
	asgn.PUT("/:id/submissions/:submissionId",
		middleware.RequireRoles("gadik", "operator"),
		h.Assignment.GradeSubmission,
	)
}
