package handler

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentController interface {
	Create(c *fiber.Ctx) error
}

type commentController struct {
	commentService services.CommentService
}

func NewCommentController(s services.CommentService) CommentController {
	return commentController{
		commentService: s,
	}
}

type CommentBody struct {
	Body   string    `validate:"required" json:"body"`
	PostID uuid.UUID `validate:"required" json:"post_id"`
}

func (cc commentController) Create(c *fiber.Ctx) error {
	comment := new(CommentBody)

	if err := c.BodyParser(comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed parse body",
			"details": err.Error(),
		})
	}

	errors := utils.ValidateStruct(comment)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed validate",
			"details": errors,
		})
	}

	jwtMeta, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed extract token metadata",
			"details": err.Error(),
		})
	}

	post, err := cc.commentService.FindPostById(comment.PostID.String())
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "error when query post",
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound && post.ID.String() == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "post not found",
		})
	}

	_, err = cc.commentService.Save(entities.Comment{
		ID:     uuid.Must(uuid.NewRandom()),
		Body:   comment.Body,
		PostID: post.ID,
		UserID: jwtMeta.UserID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "failed saving comment",
			"error":   true,
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success submit comment",
	})
}
