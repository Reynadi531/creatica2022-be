package handler

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(s services.AuthService) AuthController {
	return authController{
		authService: s,
	}
}

func (authController) Login(c *fiber.Ctx) error {
	login := new(entities.User)

	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed parse body",
			"details": err.Error(),
		})
	}

	errors := utils.ValidateStruct(login)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed validate",
			"details": errors,
		})
	}

	return c.SendString("Ok")
}

func (authController) Register(c *fiber.Ctx) error {
	panic("unimplemented")
}
