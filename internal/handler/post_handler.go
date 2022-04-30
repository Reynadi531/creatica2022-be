package handler

import (
	"strconv"

	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/database"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PostController interface {
	Create(c *fiber.Ctx) error
	ViewAll(c *fiber.Ctx) error
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

	user, err := p.postService.GetUserById(jwtMeta.UserID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed get user",
			"details": err.Error(),
		})
	}

	_, err = p.postService.Save(entities.Post{
		ID:     uuid.Must(uuid.NewRandom()),
		Body:   post.Body,
		UserID: jwtMeta.UserID,
		User:   user,
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

func (p postController) ViewAll(c *fiber.Ctx) error {
	var limit int
	var page int
	var sort string

	if c.Query("limit") == "" {
		limit = 1
	} else {
		limit, _ = strconv.Atoi(c.Query("limit"))
	}

	if c.Query("page") == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(c.Query("page"))
	}

	if c.Query("sort") == "" {
		sort = "id"
	} else {
		sort = c.Query("sort")
	}

	pageInfo, data, err := p.postService.List(database.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed get post",
			"details": err.Error(),
		})
	}

	type SendablePost struct {
		ID        uuid.UUID `json:"id"`
		Body      string    `json:"body"`
		Username  string    `json:"username"`
		CreatedAt int64     `json:"created_at"`
		UpdatedAt int64     `json:"updated_at"`
	}

	var sendablePosts []SendablePost
	for _, post := range data {
		sendablePosts = append(sendablePosts, SendablePost{
			ID:        post.ID,
			Body:      post.Body,
			Username:  post.User.Username,
			CreatedAt: post.CreatedAt.Unix(),
			UpdatedAt: post.UpdatedAt.Unix(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    fiber.StatusOK,
		"message":   "success get post",
		"page_info": pageInfo,
		"posts":     sendablePosts,
	})
}
