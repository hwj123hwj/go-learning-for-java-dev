package main

import (
	"log"
	"user-api-advanced/config"
	"user-api-advanced/controller"
	"user-api-advanced/model"
	"user-api-advanced/repository"
	"user-api-advanced/router"
	"user-api-advanced/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Config loaded successfully")

	if err := config.InitDB(cfg); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database connected successfully")

	if err := config.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(userService)

	r := router.SetupRouter(userController, authController)

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
