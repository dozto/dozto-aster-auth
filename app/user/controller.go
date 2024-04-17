package user

import (
	"errors"
	"time"

	"github.com/dozto/dozto-aster-auth/pkg"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	model *UserModel
}

func NewUserController(model *UserModel) *UserController {
	return &UserController{
		model: model,
	}
}

type CreateUserReq struct {
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"` //  TODO: validate:"phone"
	Password string `json:"password" validate:"required"`
	Meta     bson.M `json:"meta"`
}

func (u *UserController) Create(c fiber.Ctx) error {
	req := new(CreateUserReq)

	{ // Bind & Validate request body
		ErrFailedParseBody := errors.New("failed parse request body")
		if err := c.Bind().Body(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": ErrFailedParseBody.Error(),
				"error":   err,
			})
		}

		log.Info().Msgf("[user.create] create user: %s / %s ", req.Email, req.Phone)

		ErrInvalidValidation := errors.New("invalid request body validation")
		if err := pkg.Valtor.Validate(req); err != nil {
			log.Warn().Err(err).Msgf("[user.create] failed validate request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": ErrInvalidValidation.Error(),
				"error":   err,
			})
		}
	}

	{ // Check if user already exists
		ErrUserAlreadyExists := errors.New("user already exists")
		ErrFailedCheckUser := errors.New("failed check user existence")

		if req.Email != "" {
			user, err := u.model.GetByEmail(req.Email)

			if err != nil {
				log.Error().Err(err).Msg("[user.create] failed check email exist")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrFailedCheckUser.Error(),
				})
			}
			if user != nil {
				log.Warn().Msgf("[user.create] user email exist already: %s / %s", user.Email, user.ID.String())
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"message": ErrUserAlreadyExists.Error(),
				})
			}
		}

		if req.Phone != "" {
			user, err := u.model.GetByPhone(req.Phone)

			if err != nil {
				log.Error().Err(err).Msg("[user.create] failed check email exist")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrFailedCheckUser.Error(),
				})
			}
			if user != nil {
				log.Warn().Msgf("[user.create] user phone exist already: %s / %s", user.Phone, user.ID.String())
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"message": ErrUserAlreadyExists.Error(),
				})
			}
		}
	}

	// TODO: Check Email or Phone Verify Code

	{ // Hash password & Save user
		hashPassword, err := pkg.HashPassword(req.Password)
		ErrFailedHashPassword := errors.New("failed hash password")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": ErrFailedHashPassword.Error(),
				"error":   err,
			})
		}

		userDoc := &UserDoc{
			ID:             primitive.NewObjectID(),
			Email:          req.Email,
			Phone:          req.Phone,
			HashedPassword: hashPassword,
			IsSuper:        false,
			Meta:           req.Meta,
			CreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		}

		id, err := u.model.CreateUser(userDoc)
		ErrFailedSaveUser := errors.New("failed save user")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": ErrFailedSaveUser.Error(),
				"error":   err,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	}
}

// TODO:获取当前登陆用户信息

func (u *UserController) GetById(c fiber.Ctx) error {
	id := c.Params("id")

	ErrFailedGetUser := errors.New("failed to get user")
	user, err := u.model.GetUserByID(id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get user: %s", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrFailedGetUser.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	return c.JSON(user)
}
