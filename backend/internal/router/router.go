package router

import (
	"sangpalito/backend/internal/handler"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/pkg/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandlerRegistry struct {
	MasterHandler *handler.MasterHandler
}

func Setup(reg *HandlerRegistry) *gin.Engine {
	r := gin.New()

	// Middleware bawaan
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(middleware.ErrorHandler())

	// Endpoint Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.OK(gin.H{
			"status": "healthy",
			"app":    "AntriSehat API",
		}))
	})

	v1 := r.Group("/api/v1")
	{
		// Fasilitas
		v1.GET("/facilities", reg.MasterHandler.GetFacilities)

		// Poliklinik
		v1.GET("/polyclinics", reg.MasterHandler.GetPolyclinics)
		v1.POST("/polyclinics", reg.MasterHandler.CreatePolyclinic)
		v1.PATCH("/polyclinics/:id", reg.MasterHandler.UpdatePolyclinic)

		// Dokter
		v1.GET("/doctors", reg.MasterHandler.GetDoctors)
		v1.POST("/doctors", reg.MasterHandler.CreateDoctor)

		// Jadwal & Kuota
		v1.GET("/schedules", reg.MasterHandler.GetSchedules)
		v1.POST("/schedules", reg.MasterHandler.CreateSchedule)
		v1.DELETE("/schedules/:id", reg.MasterHandler.DeleteSchedule)
	}

	return r
}
