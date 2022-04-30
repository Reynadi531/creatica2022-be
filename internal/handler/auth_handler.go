package handler

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/services"
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(s services.AuthService) AuthController {
	return authController{
		authService: s,
	}
}

func (u authController) Login(c *fiber.Ctx) error {
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

	user, err := u.authService.FindByUsername(login.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "erorr when query user",
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound && user.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "user not found",
		})
	}

	if utils.ComparePassword(user.Password, login.Password) == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "wrong credentials",
		})
	}

	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "failed generating token",
			"error":   true,
			"details": err.Error(),
		})
	}

	if user.RefreshToken == "" {
		rt, _ := utils.GenerateRefreshToken()
		user.RefreshToken = rt
		_, err := u.authService.Update(user)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "failed saving token",
				"error":   true,
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": fiber.StatusOK,
			"token": fiber.Map{
				"access_token":  token,
				"refresh_token": rt,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"token": fiber.Map{
			"access_token":  token,
			"refresh_token": user.RefreshToken,
		},
	})
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

	user, err := u.authService.FindByUsername(register.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "erorr when query user",
			"details": err.Error(),
		})
	}

	if err != gorm.ErrRecordNotFound && user.Username == register.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "username already registered",
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
		ID:       uuid.Must(uuid.NewRandom()),
		Username: register.Username,
		Password: passwordHashed,
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

type RereshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (u authController) Refresh(c *fiber.Ctx) error {
	rtBody := new(RereshToken)

	if err := c.BodyParser(rtBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed parse body",
			"details": err.Error(),
		})
	}

	errors := utils.ValidateStruct(rtBody)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "failed validate",
			"details": errors,
		})
	}

	user, err := u.authService.FindByRefreshToken(rtBody.RefreshToken)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "erorr when query user",
			"details": err.Error(),
		})
	}

	if err == gorm.ErrRecordNotFound && user.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"error":   true,
			"message": "refresh token not found",
		})
	}

	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "failed generating token",
			"error":   true,
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"token": fiber.Map{
			"access_token":  token,
			"refresh_token": user.RefreshToken,
		},
	})

}
