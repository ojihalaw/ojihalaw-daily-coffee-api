package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/usecase"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Log     *logrus.Logger
	UseCase *usecase.OrderUseCase
}

func NewOrderController(useCase *usecase.OrderUseCase, logger *logrus.Logger) *OrderController {
	return &OrderController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *OrderController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateOrderRequest)

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
		JSON(utils.DefaultSuccessResponse(fiber.StatusCreated, "order created successfully"))
}

// func (c *OrderController) FindAll(ctx *fiber.Ctx) error {
// 	req := &utils.PaginationRequest{
// 		Page:    ctx.QueryInt("page", 1),
// 		Limit:   ctx.QueryInt("limit", 10),
// 		OrderBy: ctx.Query("order_by", "created_at"),
// 		SortBy:  ctx.Query("sort_by", "desc"),
// 		Search:  ctx.Query("search", ""),
// 	}

// 	categories, pagination, err := c.UseCase.FindAll(ctx.Context(), req)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).
// 			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error()))
// 	}

// 	return ctx.Status(fiber.StatusOK).
// 		JSON(utils.SuccessResponseWithPagination(fiber.StatusOK, "get list user successfully", categories, pagination))
// }

// func (c *UserController) FindByID(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")

// 	user, err := c.UseCase.FindByID(ctx.Context(), id)
// 	if err != nil {
// 		if errors.Is(err, utils.ErrNotFound) {
// 			return ctx.Status(fiber.StatusNotFound).
// 				JSON(utils.ErrorResponse(fiber.StatusNotFound, "user not found"))
// 		}

// 		return ctx.Status(fiber.StatusInternalServerError).
// 			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
// 	}

// 	return ctx.Status(fiber.StatusOK).
// 		JSON(utils.SuccessResponse(fiber.StatusOK, "get detail user successfully", user))
// }

// func (c *UserController) Update(ctx *fiber.Ctx) error {
// 	request := new(model.UpdateUserRequest)
// 	id := ctx.Params("id")

// 	err := ctx.BodyParser(request)
// 	if err != nil {
// 		c.Log.Warnf("Failed to parse request body : %+v", err)
// 		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
// 	}

// 	err = c.UseCase.Update(ctx.Context(), id, request)
// 	if err != nil {
// 		if errors.Is(err, utils.ErrNotFound) {
// 			return ctx.Status(fiber.StatusNotFound).
// 				JSON(utils.ErrorResponse(fiber.StatusNotFound, "user not found"))
// 		}

// 		return ctx.Status(fiber.StatusInternalServerError).
// 			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
// 	}

// 	return ctx.Status(fiber.StatusOK).
// 		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "update user successfully"))
// }

// func (c *UserController) Delete(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")

// 	err := c.UseCase.Delete(ctx.Context(), id)
// 	if err != nil {
// 		if errors.Is(err, utils.ErrNotFound) {
// 			return ctx.Status(fiber.StatusNotFound).
// 				JSON(utils.ErrorResponse(fiber.StatusNotFound, "user not found"))
// 		}

// 		return ctx.Status(fiber.StatusInternalServerError).
// 			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
// 	}

// 	return ctx.Status(fiber.StatusOK).
// 		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "delete user successfully"))
// }
