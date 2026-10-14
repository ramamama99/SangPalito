package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/handler"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/repository/postgres"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/internal/websocket"
	"sangpalito/backend/internal/worker"
	"sangpalito/backend/pkg/crypto"
	"sangpalito/backend/pkg/jwt"
)

// Setup initializes all routes, handlers, services, and repositories.
func Setup(
	r *gin.Engine,
	pool *pgxpool.Pool,
	jwtManager *jwt.Manager,
	nikCrypto *crypto.NIKCrypto,
	hub *websocket.Hub,
	notifWorker *worker.NotificationWorker,
	defaultFacilityID uuid.UUID,
	rateLimiter *middleware.SessionRateLimiter,
) {
	validate := validator.New()

	// ── Repositories ─────────────────────────────────────────
	userRepo := postgres.NewUserRepo(pool)
	patientRepo := postgres.NewPatientRepo(pool)
	queueRepo := postgres.NewQueueRepo(pool)
	appointmentRepo := postgres.NewAppointmentRepo(pool)
	auditLogRepo := postgres.NewAuditLogRepo(pool)
	scheduleRepo := postgres.NewScheduleRepo(pool)
	polyclinicRepo := postgres.NewPolyclinicRepo(pool)
	visitRepo := postgres.NewVisitRepo(pool)

	// ── Services ─────────────────────────────────────────────
	authService := service.NewAuthService(userRepo, auditLogRepo, jwtManager)
	patientService := service.NewPatientService(patientRepo, visitRepo, auditLogRepo, nikCrypto)
	queueService := service.NewQueueService(queueRepo, auditLogRepo, hub, notifWorker)
	publicService := service.NewPublicService(pool, patientRepo, appointmentRepo, queueRepo, scheduleRepo, polyclinicRepo, nikCrypto)
	userService := service.NewUserService(userRepo)

	// ── Handlers ─────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authService, validate, defaultFacilityID)
	patientHandler := handler.NewPatientHandler(patientService, validate)
	queueHandler := handler.NewQueueHandler(queueService, validate)
	publicHandler := handler.NewPublicHandler(publicService, validate, defaultFacilityID)
	userHandler := handler.NewUserHandler(userService, validate)

	// ── Global Middleware ────────────────────────────────────
	r.Use(middleware.Recover())
	r.Use(middleware.RequestLogger())
	if rateLimiter != nil {
		r.Use(rateLimiter.Middleware())
	}
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api/v1")

	// ── Public Routes (no auth) ─────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	public := api.Group("/public")
	{
		public.GET("/polyclinics", publicHandler.ListPolyclinics)
		public.GET("/schedules", publicHandler.ListSchedules)
		public.POST("/appointments", publicHandler.CreateAppointment)
		public.GET("/appointments/:token", publicHandler.GetAppointmentStatus)
		public.POST("/queues/checkin", publicHandler.CheckIn)
		public.GET("/queues/:id/estimate", publicHandler.GetQueueEstimate)
	}

	// ── Authenticated Routes ────────────────────────────────
	authenticated := api.Group("")
	authenticated.Use(middleware.Auth(jwtManager))
	{
		// Auth
		authenticated.POST("/auth/logout", authHandler.Logout)

		// Patients (registrar, medic)
		patients := authenticated.Group("/patients")
		{
			patients.POST("", middleware.RequireRoles("registrar"), patientHandler.Create)
			patients.GET("", middleware.RequireRoles("registrar", "medic"), patientHandler.Search)
			patients.GET("/:id", middleware.RequireRoles("registrar", "medic"), patientHandler.Get)
			patients.PATCH("/:id", middleware.RequireRoles("registrar"), patientHandler.Update)
			patients.DELETE("/:id", middleware.RequireRoles("admin"), patientHandler.Delete)
			patients.GET("/:id/visits", middleware.RequireRoles("medic"), patientHandler.Visits)
		}

		// Queues (staff)
		queues := authenticated.Group("/queues")
		{
			queues.GET("", middleware.RequireRoles("registrar", "medic"), queueHandler.List)
			queues.PATCH("/:id/status", middleware.RequireRoles("medic", "registrar"), queueHandler.UpdateStatus)
		}

		// Users (admin only)
		users := authenticated.Group("/users")
		users.Use(middleware.RequireRoles("admin"))
		{
			users.GET("", userHandler.List)
			users.POST("", userHandler.Create)
			users.PATCH("/:id", userHandler.Update)
			users.PATCH("/:id/password", userHandler.ResetPassword)
		}
	}
}
