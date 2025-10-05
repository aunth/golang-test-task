package main

import (
	"log"
	"os"

	"spy-cat-agency/internal/config"
	"spy-cat-agency/internal/database"
	"spy-cat-agency/internal/handlers"
	"spy-cat-agency/internal/middleware"
	"spy-cat-agency/internal/repository"
	"spy-cat-agency/internal/routes"
	"spy-cat-agency/internal/services"

	_ "spy-cat-agency/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	catRepo := repository.NewCatRepository(db)
	missionRepo := repository.NewMissionRepository(db)
	targetRepo := repository.NewTargetRepository(db)

	catService := services.NewCatService(catRepo)
	missionService := services.NewMissionService(missionRepo, targetRepo, catRepo)

	catHandler := handlers.NewCatHandler(catService)
	missionHandler := handlers.NewMissionHandler(missionService)

	router := gin.Default()

	router.Use(middleware.LoggingMiddleware())

	router.Use(middleware.CORSMiddleware())

	routes.SetupRoutes(router, catHandler, missionHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
