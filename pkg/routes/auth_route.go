package routes

import (
	"github.com/Reynadi531/creatica2022-be/internal/handler"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisterRouteAuth(app *fiber.App, db *gorm.DB) {
	// New Repository
	authRepository := repository.NewAuthRepository(db)

	if err := authRepository.Migrate(); err != nil {
		log.Error().Err(err).Msg("user model failed to migrate")
	}

	// New Services
	userService := services.NewAuthService(authRepository)

	// New Controller
	authController := handler.NewAuthController(userService)

	// Auth Routes
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.Refresh)
}
