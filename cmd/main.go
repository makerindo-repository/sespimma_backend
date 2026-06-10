package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"sespima_api/config"
	"sespima_api/internal/handler"
	"sespima_api/internal/repository"
	"sespima_api/internal/service"
	"sespima_api/internal/ws"
	"sespima_api/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	if err := os.MkdirAll("uploads/profile", 0755); err != nil {
		log.Fatalf("failed to create uploads dir: %v", err)
	}

	cfg := config.Load()
	config.ConnectDB(cfg)
	db := config.DB

	// ── Repositories ────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	locationLogRepo := repository.NewLocationLogRepository(db)
	kegiatanRepo := repository.NewKegiatanRepository(db)
	absensiRepo := repository.NewAbsensiRepository(db)
	izinRepo := repository.NewIzinRepository(db)
	akademikRepo := repository.NewAkademikRepository(db)
	jasmaniRepo := repository.NewJasmaniRepository(db)
	mentalRepo := repository.NewMentalRepository(db)
	sociometryRepo := repository.NewSociometryRepository(db)
	healthRepo := repository.NewHealthRepository(db)
	rewardRepo := repository.NewRewardRepository(db)
	punishmentRepo := repository.NewPunishmentRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)

	// ── Services ────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo, locationLogRepo)
	kegiatanSvc := service.NewKegiatanService(kegiatanRepo, absensiRepo, izinRepo)
	assessmentSvc := service.NewAssessmentService(akademikRepo, jasmaniRepo, mentalRepo, sociometryRepo, healthRepo)
	rewardSvc := service.NewRewardService(rewardRepo)
	punishmentSvc := service.NewPunishmentService(punishmentRepo)
	assignmentSvc := service.NewAssignmentService(assignmentRepo)

	// ── WebSocket hub ────────────────────────────────────────────────────
	hub := ws.NewHub()
	go hub.Run()

	// ── Handlers ────────────────────────────────────────────────────────
	handlers := &router.Handlers{
		Auth:       handler.NewAuthHandler(authSvc),
		User:       handler.NewUserHandler(userSvc, hub),
		Kegiatan:   handler.NewKegiatanHandler(kegiatanSvc),
		Assessment: handler.NewAssessmentHandler(assessmentSvc),
		Reward:     handler.NewRewardHandler(rewardSvc),
		Punishment: handler.NewPunishmentHandler(punishmentSvc),
		Assignment: handler.NewAssignmentHandler(assignmentSvc),
		WS:         handler.NewWSHandler(hub, cfg.JWTSecret),
		Serdik:     handler.NewSerdikHandler(),
		Dashboard:  handler.NewDashboardHandler(kegiatanSvc, rewardSvc, punishmentSvc),
		RP:         handler.NewRewardPunishmentHandler(),
	}

	// ── HTTP Server ──────────────────────────────────────────────────────
	r := gin.Default()
	router.Register(r, handlers, cfg.JWTSecret)

	log.Printf("server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
