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
	ViewDetail(c *fiber.Ctx) error
	CreateReplies(c *fiber.Ctx) error
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

	commentData, err := cc.commentService.Save(entities.Comment{
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
		"status":     fiber.StatusOK,
		"message":    "success submit comment",
		"comment_id": commentData.ID.String(),
	})
}

type ReplayBody struct {
	Body string `validate:"required" json:"body"`
}

func (cc commentController) CreateReplies(c *fiber.Ctx) error {
	replayBody := new(ReplayBody)
	commentId := c.Params("id")

	if err := c.BodyParser(replayBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed parse body",
			"details": err.Error(),
		})
	}

	errors := utils.ValidateStruct(replayBody)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed validate",
			"details": errors,
		})
	}

	comment, err := cc.commentService.FindCommentById(commentId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "error when query comment",
			"error":   true,
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound && comment.ID.String() == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "comment not found",
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

	reply, err := cc.commentService.SaveReply(entities.Reply{
		ID:        uuid.Must(uuid.NewRandom()),
		Body:      replayBody.Body,
		CommentID: comment.ID,
		UserID:    jwtMeta.UserID,
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
		"status":    fiber.StatusOK,
		"message":   "success submit reply",
		"replay_id": reply.ID,
	})
}

func (cc commentController) ViewDetail(c *fiber.Ctx) error {
	commentId := c.Params("id")

	comment, err := cc.commentService.FindCommentById(commentId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "error when query comment",
			"error":   true,
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound && comment.ID.String() == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "comment not found",
			"error":   true,
		})
	}

	reply, err := cc.commentService.FindReplyByCommentId(comment.ID.String())
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "error when query reply",
			"error":   true,
			"details": err.Error(),
		})
	}

	type SendableComment struct {
		ID        uuid.UUID `json:"id"`
		Body      string    `json:"body"`
		Username  string    `json:"username"`
		CreatedAt int64     `json:"created_at"`
		UpdatedAt int64     `json:"updated_at"`
	}

	sendableComment := SendableComment{
		ID:        comment.ID,
		Body:      comment.Body,
		Username:  comment.User.Username,
		CreatedAt: comment.CreatedAt.Unix(),
		UpdatedAt: comment.UpdatedAt.Unix(),
	}

	type SendableReply struct {
		ID        uuid.UUID `json:"id"`
		Body      string    `json:"body"`
		Username  string    `json:"username"`
		CreatedAt int64     `json:"created_at"`
		UpdatedAt int64     `json:"updated_at"`
	}
	var sendableReply []SendableReply

	for _, i := range reply {
		sendableReply = append(sendableReply, SendableReply{
			ID:        i.ID,
			Body:      i.Body,
			Username:  i.User.Username,
			CreatedAt: i.CreatedAt.Unix(),
			UpdatedAt: i.UpdatedAt.Unix(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success get comment",
		"comment": sendableComment,
		"reply":   sendableReply,
	})
}
