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

func RegisterRouteComment(app *fiber.App, db *gorm.DB) {
	//	New Repository
	commentRepository := repository.NewCommentRepository(db)

	if err := commentRepository.Migrate(); err != nil {
		log.Error().Err(err).Msg("comment model failed to migrate")
	}

	//	New Services
	commentService := services.CommentService(commentRepository)

	//	New Controller
	commentController := handler.NewCommentController(commentService)

	//	Comment Routes
	comment := app.Group("/api/v1/comment")
	comment.Post("/", middleware.JWTProtected(), commentController.Create)
	comment.Get("/:id", middleware.JWTProtected(), commentController.ViewDetail)
	comment.Post("/:id/reply", middleware.JWTProtected(), commentController.CreateReplies)
}
