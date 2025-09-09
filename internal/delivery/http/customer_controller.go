package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/usecase"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
)

type CustomerController struct {
	Log     *logrus.Logger
	UseCase *usecase.CustomerUseCase
}

func NewCustomerController(useCase *usecase.CustomerUseCase, logger *logrus.Logger) *CustomerController {
	return &CustomerController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *CustomerController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterCustomerRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	err = c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)

		switch {
		case errors.Is(err, utils.ErrValidation):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))

		case errors.Is(err, utils.ErrConflict):
			return ctx.Status(fiber.StatusConflict).
				JSON(utils.ErrorResponse(fiber.StatusConflict, err.Error()))

		default: // internal error
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
		}
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(utils.SuccessResponse(fiber.StatusCreated, "customer registered successfully", nil))
}

func (c *CustomerController) FindAll(ctx *fiber.Ctx) error {
	req := &utils.PaginationRequest{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", 10),
		OrderBy: ctx.Query("order_by", "created_at"),
		SortBy:  ctx.Query("sort_by", "desc"),
		Search:  ctx.Query("search", ""),
	}

	categories, pagination, err := c.UseCase.FindAll(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponseWithPagination(fiber.StatusOK, "get list customer successfully", categories, pagination))
}

func (c *CustomerController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := c.UseCase.FindByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "customer not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "get detail customer successfully", user))
}

func (c *CustomerController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateCustomerRequest)
	id := ctx.Params("id")

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	err = c.UseCase.Update(ctx.Context(), id, request)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "customer not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "update customer successfully"))
}

func (c *CustomerController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.UseCase.Delete(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "customer not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "delete customer successfully"))
}
