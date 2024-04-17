package auth

import (
	"errors"

	"github.com/dozto/dozto-aster-auth/app/user"
	"github.com/dozto/dozto-aster-auth/pkg"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	userModel *user.UserModel
}

func NewAuthController(model *user.UserModel) *AuthController {
	return &AuthController{
		userModel: model,
	}
}

type LoginType string

const (
	Email  LoginType = "email"
	Phone  LoginType = "phone"
	WeChat LoginType = "wechat"
)

type LoginRequest struct {
	Type     LoginType `json:"type" validate:"required"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password" validate:"required"`
}

func (a *AuthController) Login(c fiber.Ctx) error {
	req := new(LoginRequest)

	{ // Bind & Validate request body
		ErrFailedParseBody := errors.New("failed parse request body")
		if err := c.Bind().Body(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": ErrFailedParseBody.Error(),
				"error":   err,
			})
		}

		log.Info().Msgf("[auth.login] login user: %s / %s ", req.Email, req.Phone)

		ErrInvalidValidation := errors.New("invalid request body validation")
		if err := pkg.Valtor.Validate(req); err != nil {
			log.Warn().Err(err).Msgf("[auth.login] failed validate request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": ErrInvalidValidation.Error(),
				"error":   pkg.Valtor.TransErrorToResponse(&err),
			})
		}
	}

	var user *user.UserDoc
	{ // Get target login user
		switch req.Type {
		case Email:
			user, _ = a.userModel.GetByEmail(req.Email)
		case Phone:
			user, _ = a.userModel.GetByPhone(req.Phone)
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid login method",
			})
		}
	}

	{ // Check password correct and target user exist.
		if user != nil {
			if pkg.ValidateHashedPassword(req.Password, user.HashedPassword) {
				// TODO: Generate JWT token and return
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"message": "login success",
				})
			}
		}

		// Failed validate password or user not exist
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "incorrect credential or user not exist",
		})
	}
}
