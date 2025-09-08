package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/usecase"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	Log     *logrus.Logger
	UseCase *usecase.ProductUseCase
}

func NewProductController(useCase *usecase.ProductUseCase, logger *logrus.Logger) *ProductController {
	return &ProductController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ProductController) Create(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name")
	slug := ctx.FormValue("slug")
	sku := ctx.FormValue("sku")
	variant := ctx.FormValue("variant")
	price, _ := strconv.Atoi(ctx.FormValue("price"))
	stock, _ := strconv.Atoi(ctx.FormValue("stock"))
	description := ctx.FormValue("description")
	categoryIDStr := ctx.FormValue("category_id")

	// Convert categoryID ke UUID
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid category_id"))
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Image is required"))
	}

	request := &model.CreateProductRequest{
		Name:        name,
		Slug:        slug,
		SKU:         sku,
		Variant:     variant,
		Price:       price,
		Stock:       stock,
		Description: description,
		CategoryID:  categoryID,
	}

	err = c.UseCase.Create(ctx.UserContext(), request, file)
	if err != nil {
		c.Log.Warnf("Failed to create product : %+v", err)

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
		JSON(utils.DefaultSuccessResponse(fiber.StatusCreated, "product created successfully"))
}

func (c *ProductController) FindAll(ctx *fiber.Ctx) error {
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
		JSON(utils.SuccessResponseWithPagination(fiber.StatusOK, "get list product successfully", categories, pagination))
}

func (c *ProductController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	product, err := c.UseCase.FindByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "product not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "get detail product successfully", product))
}

func (c *ProductController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Ambil form-data
	name := ctx.FormValue("name")
	slug := ctx.FormValue("slug")
	sku := ctx.FormValue("sku")
	variant := ctx.FormValue("variant")
	priceStr := ctx.FormValue("price")
	stockStr := ctx.FormValue("stock")
	description := ctx.FormValue("description")
	categoryIDStr := ctx.FormValue("category_id")

	var price, stock int
	if priceStr != "" {
		price, _ = strconv.Atoi(priceStr)
	}
	if stockStr != "" {
		stock, _ = strconv.Atoi(stockStr)
	}

	var categoryID uuid.UUID
	if categoryIDStr != "" {
		parsed, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Invalid category_id"))
		}
		categoryID = parsed
	}

	// File opsional
	file, _ := ctx.FormFile("image")

	// Build request
	request := &model.UpdateProductRequest{
		Name:        name,
		Slug:        slug,
		SKU:         sku,
		Variant:     variant,
		Price:       price,
		Stock:       stock,
		Description: description,
		CategoryID:  categoryID,
	}

	err := c.UseCase.Update(ctx.Context(), id, request, file)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "product not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "update product successfully"))
}

func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.UseCase.Delete(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "product not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "delete product successfully"))
}

func (c *ProductController) FindSpecialProduct(ctx *fiber.Ctx) error {
	product, err := c.UseCase.FindSpecialProduct(ctx.Context())
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))

		case errors.Is(err, utils.ErrNotFound):
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, err.Error()))

		default: // internal error
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
		}
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "get special product successfully", product))
}

func (c *ProductController) UpdateSpecialProduct(ctx *fiber.Ctx) error {
	request := new(model.UpdateSpecialProductRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	err = c.UseCase.SetSpecialProduct(ctx.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrValidation):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))

		default: // internal error
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
		}
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "update special product successfully"))
}
