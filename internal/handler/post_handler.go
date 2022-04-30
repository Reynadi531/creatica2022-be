package handler

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PostController interface {
	Create(c *fiber.Ctx) error
}

type postController struct {
	postService services.PostService
}

func NewPostController(s services.PostService) PostController {
	return postController{
		postService: s,
	}
}

type PostBody struct {
	Body string `validate:"required,min=3" json:"body"`
}

func (p postController) Create(c *fiber.Ctx) error {
	post := new(PostBody)

	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed parse body",
			"details": err.Error(),
		})
	}

	errors := utils.ValidateStruct(post)
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

	_, err = p.postService.Save(entities.Post{
		ID:     uuid.Must(uuid.NewRandom()),
		Body:   post.Body,
		UserID: jwtMeta.UserID,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed save post",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "success create post",
	})
}
