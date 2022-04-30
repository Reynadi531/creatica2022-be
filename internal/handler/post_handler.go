package handler

import (
	"strconv"

	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/database"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostController interface {
	Create(c *fiber.Ctx) error
	ViewAll(c *fiber.Ctx) error
	ViewDetail(c *fiber.Ctx) error
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

	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed get post",
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"error":   true,
			"message": "post not found",
		})
	}

	type SendablePost struct {
		ID           uuid.UUID `json:"id"`
		Body         string    `json:"body"`
		CommentCount int64     `json:"comment_count"`
		Username     string    `json:"username"`
		CreatedAt    int64     `json:"created_at"`
		UpdatedAt    int64     `json:"updated_at"`
	}

	var sendablePosts []SendablePost
	for _, post := range data {
		count, err := p.postService.CountCommentOnPost(post.ID.String())
		if err != nil && err != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "failed counting comment",
				"error":   true,
				"details": err.Error(),
			})
		}
		sendablePosts = append(sendablePosts, SendablePost{
			ID:           post.ID,
			Body:         post.Body,
			CommentCount: count,
			Username:     post.User.Username,
			CreatedAt:    post.CreatedAt.Unix(),
			UpdatedAt:    post.UpdatedAt.Unix(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    fiber.StatusOK,
		"message":   "success get post",
		"page_info": pageInfo,
		"posts":     sendablePosts,
	})
}

func (p postController) ViewDetail(c *fiber.Ctx) error {
	postID := c.Params("id")

	post, err := p.postService.GetPostById(postID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed get post",
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"error":   true,
			"message": "post not found",
		})
	}

	comment, err := p.postService.GetCommentByPostId(post.ID.String())
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"error":   true,
			"message": "error when query comment",
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

	var sendablePost SendablePost
	sendablePost = SendablePost{
		ID:        post.ID,
		Body:      post.Body,
		Username:  post.User.Username,
		CreatedAt: post.CreatedAt.Unix(),
		UpdatedAt: post.UpdatedAt.Unix(),
	}

	type SendableComment struct {
		ID        uuid.UUID `json:"id"`
		Body      string    `json:"body"`
		Username  string    `json:"username"`
		CreatedAt int64     `json:"created_at"`
		UpdatedAt int64     `json:"updated_at"`
	}
	var sendableComment []SendableComment
	for _, i := range comment {
		sendableComment = append(sendableComment, SendableComment{
			ID:        i.ID,
			Body:      i.Body,
			Username:  i.User.Username,
			CreatedAt: i.CreatedAt.Unix(),
			UpdatedAt: i.CreatedAt.Unix(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"post":    sendablePost,
		"comment": sendableComment,
	})

}
