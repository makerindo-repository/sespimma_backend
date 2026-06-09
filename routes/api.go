package routes

import (
	"sespima_api/controllers"
	"sespima_api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	api.POST("/login", controllers.Login)
	api.POST("/reset_password", controllers.ResetPassword)
	api.POST("/verify_reset_token", controllers.VerifyResetToken)

	api.Use(middleware.Auth())
	api.POST("/generate_reset_token", controllers.GenerateResetToken)
	api.POST("/update_password", controllers.UpdatePassword)

	userRewardCtrl := controllers.UserRewardController{DB: db}
	api.GET("/user_rewards", userRewardCtrl.GetAll)
	api.GET("/user_rewards/pending", userRewardCtrl.GetPendingApproval)
	api.POST("/user_rewards", userRewardCtrl.Create)
	api.PUT("/user_rewards/:id", userRewardCtrl.Update)
	api.DELETE("/user_rewards/:id", userRewardCtrl.Delete)
	api.PUT("/user_rewards/:id/approve", userRewardCtrl.Approve)

	punishmentLogCtrl := controllers.PunishmentLogController{DB: db}
	api.GET("/punishment_logs", punishmentLogCtrl.GetAll)
	api.GET("/punishment_logs/user/:user_id", punishmentLogCtrl.GetByUserID)
	api.POST("/punishment_logs", punishmentLogCtrl.Create)
	api.PUT("/punishment_logs/:id", punishmentLogCtrl.Update)
	api.DELETE("/punishment_logs/:id", punishmentLogCtrl.Delete)

	kegiatanCtrl := controllers.KegiatanController{DB: db}
	api.GET("/kegiatan", kegiatanCtrl.GetAll)
	api.GET("/kegiatan/rutin", kegiatanCtrl.GetRutin)
	api.GET("/kegiatan/today", kegiatanCtrl.GetToday)
	api.POST("/kegiatan", kegiatanCtrl.Create)
	api.PUT("/kegiatan/:id", kegiatanCtrl.Update)
	api.DELETE("/kegiatan/:id", kegiatanCtrl.Delete)

	user := r.Group("/users")
	{
		user.GET("", controllers.GetAllUsers)
		user.GET("/locations", controllers.GetAllUserLocations)
		user.GET("/:id/location", controllers.GetUserLocation)
		user.PUT("/:id/location", controllers.UpdateUserLocation)
		user.PUT("/me/location", middleware.Auth(), controllers.UpdateMyLocation)
	}

	akademikController := controllers.NewAkademikController()
	akademik := api.Group("/akademik")
	{
		akademik.GET("", akademikController.GetAll)
		akademik.GET("/serdik/:serdikId", akademikController.GetBySerdikID)
		akademik.POST("", akademikController.Create)
		akademik.PUT("/:id", akademikController.Update)
		akademik.DELETE("/:id", akademikController.Delete)
	}
}
