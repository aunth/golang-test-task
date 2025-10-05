package routes

import (
	"spy-cat-agency/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, catHandler *handlers.CatHandler, missionHandler *handlers.MissionHandler) {
	v1 := router.Group("/api/v1")
	{
		SetupCatRoutes(v1, catHandler)
		SetupMissionRoutes(v1, missionHandler)
		SetupTargetRoutes(v1, missionHandler)
	}
}
