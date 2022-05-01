package routes

import (
	"github.com/Reynadi531/creatica2022-be/internal/handler"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RegisterRoutePost(app *fiber.App, db *gorm.DB) {
	// New Repository
	postRepository := repository.NewPostRepository(db)

	if err := postRepository.Migrate(); err != nil {
		log.Error().Err(err).Msg("post model failed to migrate")
	}

	// New Services
	postService := services.NewPostService(postRepository)

	// New Controller
	postController := handler.NewPostController(postService)

	// Auth Routes
	post := app.Group("/api/v1/post")
	post.Post("/", middleware.JWTProtected(), postController.Create)
	post.Get("/", middleware.JWTProtected(), postController.ViewAll)
	post.Get("/me", middleware.JWTProtected(), postController.ViewSelf)
	post.Get("/:id", middleware.JWTProtected(), postController.ViewDetail)
}
