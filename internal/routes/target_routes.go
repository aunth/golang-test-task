package routes

import (
	"spy-cat-agency/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupTargetRoutes(router *gin.RouterGroup, missionHandler *handlers.MissionHandler) {

	targets := router.Group("/missions/:id/targets")
	{
		targets.POST("", missionHandler.AddTarget)
		targets.PUT("/:targetId", missionHandler.UpdateTarget)
		targets.DELETE("/:targetId", missionHandler.DeleteTarget)
		targets.PUT("/:targetId/complete", missionHandler.CompleteTarget)
		targets.PUT("/:targetId/notes", missionHandler.UpdateTargetNotes)
	}
}
