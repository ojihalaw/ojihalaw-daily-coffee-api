package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/usecase"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log     *logrus.Logger
	UseCase *usecase.AuthUseCase
}

func NewAuthController(useCase *usecase.AuthUseCase, logger *logrus.Logger) *AuthController {
	return &AuthController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (a *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		a.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	token, err := a.UseCase.Login(ctx.UserContext(), request, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		a.Log.Warnf("Failed to login : %+v", err)

		switch {
		case errors.Is(err, utils.ErrValidation):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))
		case errors.Is(err, utils.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(utils.ErrorResponse(fiber.StatusUnauthorized, err.Error()))
		case errors.Is(err, utils.ErrInvalidPassword):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))
		case errors.Is(err, utils.ErrInvalidEmail):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))

		default:
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
		}
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "Login successfully", token))
}

func (a *AuthController) LoginClientWithEmail(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		a.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	token, err := a.UseCase.ClientLoginWithEmail(ctx.UserContext(), request, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		a.Log.Warnf("Failed to login : %+v", err)

		switch {
		case errors.Is(err, utils.ErrValidation):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))
		case errors.Is(err, utils.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(utils.ErrorResponse(fiber.StatusUnauthorized, err.Error()))
		case errors.Is(err, utils.ErrInvalidPassword):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))
		case errors.Is(err, utils.ErrInvalidEmail):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))

		default:
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
		}
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "Login successfully", token))
}
