package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model/converter"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validator          *utils.Validator
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	categoryRepository *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Log:                logger,
		Validator:          validator,
		CategoryRepository: categoryRepository,
	}
}

func (c *CategoryUseCase) Create(ctx context.Context, request *model.CreateCategoryRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(c.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	// check duplicate
	exists, err := c.CategoryRepository.ExistsByName(c.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "category name already exist")
	}

	slug := utils.GenerateSlug(request.Name)

	category := &entity.Category{
		Name: request.Name,
		Slug: slug,
	}

	if err := c.CategoryRepository.Create(c.DB.WithContext(ctx), category); err != nil {
		c.Log.Warnf("Failed create category to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *CategoryUseCase) FindAll(ctx context.Context, pagination *utils.PaginationRequest) ([]model.CategoryResponse, *utils.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var categories []entity.Category

	total, err := c.CategoryRepository.FindAll(c.DB.WithContext(ctx), &categories, pagination)
	if err != nil {
		c.Log.Warnf("Failed find all category from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	responses := make([]model.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *converter.CategoryToResponse(&category)
	}

	totalPage := int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginationRes := &utils.PaginationResponse{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		OrderBy:   pagination.OrderBy,
		SortBy:    pagination.SortBy,
		Search:    pagination.Search,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses, paginationRes, nil
}

func (c *CategoryUseCase) FindByID(ctx context.Context, categoryID string) (*model.CategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var categories *entity.Category

	category, err := c.CategoryRepository.FindById(c.DB.WithContext(ctx), categories, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Category not found, id=%s", categoryID)
			return nil, utils.ErrNotFound
		}
		c.Log.Warnf("Failed find category from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return converter.CategoryToResponse(category), nil
}

func (c *CategoryUseCase) Update(ctx context.Context, categoryID string, request *model.UpdateCategoryRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	category := &entity.Category{}
	_, err := c.CategoryRepository.FindById(c.DB.WithContext(ctx), category, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Category not found, id=%s", categoryID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find category from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(c.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	// check duplicate
	exists, err := c.CategoryRepository.ExistsByName(c.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "category name already exist")
	}

	slug := utils.GenerateSlug(request.Name)

	category.Name = request.Name
	category.Slug = slug

	err = c.CategoryRepository.Update(c.DB.WithContext(ctx), category)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Category not found, id=%s", categoryID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find category from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *CategoryUseCase) Delete(ctx context.Context, categoryID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	category := &entity.Category{}
	_, err := c.CategoryRepository.FindById(c.DB.WithContext(ctx), category, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Category not found, id=%s", categoryID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find category from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.CategoryRepository.Delete(c.DB.WithContext(ctx), category)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Category not found, id=%s", categoryID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find category from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
