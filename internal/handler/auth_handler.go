package handler

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

func (u authController) Register(c *fiber.Ctx) error {
	register := new(entities.User)

	if err := c.BodyParser(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed parse body",
			"details": err.Error(),
		})
	}

	errors := utils.ValidateStruct(register)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed validate",
			"details": errors,
		})
	}

	user, err := u.authService.FindByEmail(register.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "erorr when query user",
			"details": err.Error(),
		})
	}

	if err != gorm.ErrRecordNotFound && user.Email == register.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "email already registered",
		})
	}

	passwordHashed, err := utils.GeneratePassword(register.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"error":   true,
			"message": "failed generating password",
			"details": err.Error(),
		})
	}

	_, err = u.authService.Save(entities.User{
		Email:    register.Email,
		Password: passwordHashed,
		Username: register.Username,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"error":   true,
			"message": "failed insert user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "success register",
	})

}
