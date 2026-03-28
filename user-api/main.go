package main

import (
	"log"
	"user-api/config"
	"user-api/controller"
	"user-api/repository"
	"user-api/router"
	"user-api/service"
)

func main() {
	cfg := config.LoadConfig()

	if err := config.InitDB(cfg); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database connected and migrated successfully")

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := router.SetupRouter(userController)

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
