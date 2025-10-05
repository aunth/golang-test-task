package routes

import (
	"spy-cat-agency/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupMissionRoutes(router *gin.RouterGroup, missionHandler *handlers.MissionHandler) {
	missions := router.Group("/missions")
	{
		missions.POST("", missionHandler.CreateMission)
		missions.GET("", missionHandler.ListMissions)
		missions.GET("/:id", missionHandler.GetMission)
		missions.PUT("/:id", missionHandler.UpdateMission)
		missions.DELETE("/:id", missionHandler.DeleteMission)
		missions.PUT("/:id/assign", missionHandler.AssignCat)
		missions.PUT("/:id/complete", missionHandler.CompleteMission)
	}
}
