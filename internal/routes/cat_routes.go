package routes

import (
	"spy-cat-agency/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCatRoutes(router *gin.RouterGroup, catHandler *handlers.CatHandler) {
	cats := router.Group("/cats")
	{
		cats.POST("", catHandler.CreateCat)
		cats.GET("", catHandler.ListCats)
		cats.GET("/:id", catHandler.GetCat)
		cats.PUT("/:id", catHandler.UpdateCat)
		cats.DELETE("/:id", catHandler.DeleteCat)
	}
}
