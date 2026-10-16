package router

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/handler"
	"sangpalito/backend/internal/hub"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/repository/postgres"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/internal/worker"
	"sangpalito/backend/pkg/crypto"
	"sangpalito/backend/pkg/jwt"
	"sangpalito/backend/pkg/response"
)

// Setup initializes all routes, handlers, services, and repositories.
func Setup(
	r *gin.Engine,
	pool *pgxpool.Pool,
	jwtManager *jwt.Manager,
	nikCrypto *crypto.NIKCrypto,
	hub *hub.Hub,
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
	doctorRepo := postgres.NewDoctorRepo(pool)
	facilityRepo := postgres.NewFacilityRepo(pool)
	visitRepo := postgres.NewVisitRepo(pool)
	preRepo := postgres.NewPreAnamnesisRepo(pool)
	icd10Repo := postgres.NewICD10Repository(pool)
	medicationRepo := postgres.NewMedicationRepo(pool)
	prescriptionRepo := postgres.NewPrescriptionRepo(pool)

	// Auto-seed Puskesmas common ICD-10 codes in background
	go func() {
		_ = icd10Repo.SeedCommonCodes(context.Background())
	}()

	// ── Services ─────────────────────────────────────────────
	authService := service.NewAuthService(userRepo, doctorRepo, auditLogRepo, jwtManager)
	patientService := service.NewPatientService(patientRepo, visitRepo, auditLogRepo, nikCrypto)
	queueService := service.NewQueueService(queueRepo, auditLogRepo, hub, notifWorker)
	publicService := service.NewPublicService(pool, patientRepo, appointmentRepo, queueRepo, scheduleRepo, polyclinicRepo, nikCrypto)
	userService := service.NewUserService(userRepo)
	masterService := service.NewMasterService(facilityRepo, polyclinicRepo, doctorRepo, scheduleRepo, userRepo)
	apptService := service.NewAppointmentService(appointmentRepo, preRepo, polyclinicRepo, facilityRepo)
	visitService := service.NewVisitService(visitRepo, icd10Repo, prescriptionRepo, auditLogRepo)
	prescriptionService := service.NewPrescriptionService(prescriptionRepo, medicationRepo)
	analyticsService := service.NewAnalyticsService(pool)

	// ── Handlers ─────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authService, validate, defaultFacilityID)
	patientHandler := handler.NewPatientHandler(patientService, validate)
	queueHandler := handler.NewQueueHandler(queueService, validate)
	publicHandler := handler.NewPublicHandler(publicService, validate, defaultFacilityID)
	userHandler := handler.NewUserHandler(userService, validate)
	masterHandler := handler.NewMasterHandler(masterService)
	apptHandler := handler.NewAppointmentHandler(apptService, masterService)
	visitHandler := handler.NewVisitHandler(visitService)
	prescriptionHandler := handler.NewPrescriptionHandler(prescriptionService, publicService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService, defaultFacilityID)
	wsHandler := handler.NewWSHandler(hub, pool)
	auditLogHandler := handler.NewAuditLogHandler(auditLogRepo)

	// ── Global Middleware ────────────────────────────────────
	r.Use(middleware.Recover())
	r.Use(cors.Default())
	r.Use(middleware.RequestLogger())
	if rateLimiter != nil {
		r.Use(rateLimiter.Middleware())
	}
	r.Use(middleware.ErrorHandler())

	// ── Health Check ─────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.OK(gin.H{
			"status": "healthy",
			"app":    "SangPalito AntriSehat API",
		}))
	})

	// ── WebSocket (Real-time Queues Stream) ───────────────────
	r.GET("/ws/queues", wsHandler.ServeWS)
	r.GET("/ws/queues/:id", wsHandler.ServeWS)

	api := r.Group("/api/v1")

	// ── Master Data Routes (Public / Read) ────────────────────
	api.GET("/facilities", masterHandler.GetFacilities)
	api.GET("/polyclinics", masterHandler.GetPolyclinics)
	api.GET("/doctors", masterHandler.GetDoctors)
	api.GET("/schedules", masterHandler.GetSchedules)
	api.POST("/schedules/assign", masterHandler.AssignShift)
	api.POST("/schedules/unassign", masterHandler.UnassignShift)
	api.DELETE("/schedules/:id", masterHandler.DeleteSchedule)
	api.GET("/icd10", visitHandler.SearchICD10)
	api.GET("/medications", prescriptionHandler.ListMedications)
	api.POST("/medications", prescriptionHandler.CreateMedication)
	api.PUT("/medications/:id", prescriptionHandler.UpdateMedication)
	api.DELETE("/medications/:id", prescriptionHandler.DeleteMedication)

	// ── Public Routes (no auth) ─────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	public := api.Group("/public")
	{
		public.GET("/polyclinics", publicHandler.ListPolyclinics)
		public.GET("/schedules", publicHandler.ListSchedules)
		public.POST("/appointments", publicHandler.CreateAppointment)
		public.GET("/appointments/:token", publicHandler.GetAppointmentStatus)
		public.GET("/appointments/:token/prescription", prescriptionHandler.GetByAppointmentToken)
		public.GET("/prescriptions/:prescription_no", prescriptionHandler.GetByNo)
		public.POST("/queues/checkin", publicHandler.CheckIn)
		public.GET("/queues/demo", publicHandler.GetDemoQueues)
		public.POST("/queues/reset-demo", publicHandler.ResetDemoQueue)
		public.GET("/queues/:id/estimate", publicHandler.GetQueueEstimate)
		public.GET("/queues/:id/display", wsHandler.GetPolyclinicDisplay)
	}

	// ── Authenticated Routes ────────────────────────────────
	authenticated := api.Group("")
	authenticated.Use(middleware.Auth(jwtManager))
	{
		// Auth
		authenticated.POST("/auth/logout", authHandler.Logout)
		authenticated.POST("/auth/first-login", authHandler.CompleteFirstLogin)

		// Master Data Management (Admin / Staff)
		master := authenticated.Group("")
		master.Use(middleware.RequireRoles("admin", "head"))
		{
			master.POST("/polyclinics", masterHandler.CreatePolyclinic)
			master.PATCH("/polyclinics/:id", masterHandler.UpdatePolyclinic)
			master.POST("/doctors", masterHandler.CreateDoctor)
			master.PUT("/doctors/:id", masterHandler.UpdateDoctor)
			master.DELETE("/doctors/:id", masterHandler.DeleteDoctor)
			master.POST("/schedules", masterHandler.CreateSchedule)
		}

		// Appointments & Pre-anamnesis (Registrar, Medic)
		appointments := authenticated.Group("/appointments")
		{
			appointments.GET("", apptHandler.GetAppointments)
			appointments.POST("", apptHandler.CreateWalkIn)
			appointments.DELETE("/:id", apptHandler.CancelAppointment)
			appointments.GET("/:id/pre-anamnesis", apptHandler.GetPreAnamnesis)
			appointments.PATCH("/:id/pre-anamnesis", apptHandler.UpdatePreAnamnesis)
		}

		// Patients (registrar, medic)
		patients := authenticated.Group("/patients")
		{
			patients.POST("", middleware.RequireRoles("registrar"), patientHandler.Create)
			patients.GET("", middleware.RequireRoles("registrar", "medic"), patientHandler.Search)
			patients.GET("/:id", middleware.RequireRoles("registrar", "medic"), patientHandler.Get)
			patients.PATCH("/:id", middleware.RequireRoles("registrar"), patientHandler.Update)
			patients.DELETE("/:id", middleware.RequireRoles("admin"), patientHandler.Delete)
			patients.GET("/:id/visits", middleware.RequireRoles("medic", "registrar"), patientHandler.Visits)
			patients.GET("/:id/medical-records", middleware.RequireRoles("medic", "registrar"), patientHandler.Visits)
		}

		// Queues (staff)
		queues := authenticated.Group("/queues")
		{
			queues.GET("", middleware.RequireRoles("registrar", "medic"), queueHandler.List)
			queues.PATCH("/:id/status", middleware.RequireRoles("medic", "registrar"), queueHandler.UpdateStatus)
		}

		// Clinical Visits / Rekam Medis SOAP (medic, head, admin)
		visits := authenticated.Group("/visits")
		{
			visits.POST("", middleware.RequireRoles("medic", "head", "admin"), visitHandler.CreateVisit)
			visits.POST("/audit-log", middleware.RequireRoles("medic", "head", "admin"), visitHandler.LogVisitAudit)
			visits.GET("/:id", middleware.RequireRoles("medic", "head", "admin"), visitHandler.GetVisitByID)
			visits.PATCH("/:id", middleware.RequireRoles("medic", "head", "admin"), visitHandler.UpdateVisit)
		}

		// Prescriptions / e-Resep Obat (medic, registrar, admin)
		prescriptions := authenticated.Group("/prescriptions")
		{
			prescriptions.POST("", middleware.RequireRoles("medic", "head", "admin"), prescriptionHandler.Create)
			prescriptions.PATCH("/:id/dispense", middleware.RequireRoles("medic", "registrar", "admin"), prescriptionHandler.Dispense)
		}

		// Users (admin only)
		users := authenticated.Group("/users")
		users.Use(middleware.RequireRoles("admin"))
		{
			users.GET("", userHandler.List)
			users.POST("", userHandler.Create)
			users.PATCH("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.PATCH("/:id/password", userHandler.ResetPassword)
		}

		// Analytics Dashboard (head, admin)
		analytics := authenticated.Group("/analytics")
		analytics.Use(middleware.RequireRoles("head", "admin"))
		{
			analytics.GET("/summary", analyticsHandler.GetSummary)
			analytics.GET("/busy-hours", analyticsHandler.GetBusyHours)
			analytics.GET("/top-diagnoses", analyticsHandler.GetTopDiagnoses)
			analytics.GET("/weekly-patients", analyticsHandler.GetWeeklyPatients)
			analytics.GET("/doctor-performance", analyticsHandler.GetDoctorPerformance)
		}

		// Audit Logs (admin, head - F-12)
		authenticated.GET("/audit-logs", middleware.RequireRoles("admin", "head"), auditLogHandler.List)
	}
}
