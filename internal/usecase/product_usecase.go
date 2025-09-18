package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model/converter"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validator         *utils.Validator
	Cloudinary        *cloudinary.Cloudinary
	ProductRepository *repository.ProductRepository
}

func NewProductUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator, cloudinary *cloudinary.Cloudinary, productRepository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		DB:                db,
		Log:               logger,
		Validator:         validator,
		Cloudinary:        cloudinary,
		ProductRepository: productRepository,
	}
}

func (p *ProductUseCase) Create(ctx context.Context, request *model.CreateProductRequest, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(p.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	// check duplicate
	exists, err := p.ProductRepository.ExistsByName(p.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "product name already exist")
	}

	imageURL, err := utils.UploadImageIntoCloudinary(
		p.Cloudinary,
		ctx,
		file,
		"daily-coffee/products",
	)

	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	slug := utils.GenerateSlug(request.Name)
	sku := utils.GenerateSKU(request.Name)

	product := &entity.Product{
		Name:        request.Name,
		SKU:         sku,
		Variant:     request.Variant,
		Price:       request.Price,
		Description: request.Description,
		Stock:       request.Stock,
		Slug:        slug,
		ImageURL:    imageURL,
		CategoryID:  request.CategoryID,
	}

	if err := p.ProductRepository.Create(p.DB.WithContext(ctx), product); err != nil {
		p.Log.Warnf("Failed create product to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (p *ProductUseCase) FindAll(ctx context.Context, pagination *utils.PaginationRequest) ([]model.ProductResponse, *utils.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var products []entity.Product

	total, err := p.ProductRepository.FindAllWithRedis(p.DB.WithContext(ctx), &products, pagination)
	if err != nil {
		p.Log.Warnf("Failed find all category from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	responses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *converter.ProductToResponse(&product)
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

func (p *ProductUseCase) FindByID(ctx context.Context, productID string) (*model.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var products *entity.Product

	product, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), products, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return nil, utils.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return converter.ProductToResponse(product), nil
}

func (p *ProductUseCase) Update(ctx context.Context, productID string, request *model.UpdateProductRequest, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Cari product lama
	product := &entity.Product{}
	_, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), product, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	if file != nil {
		imageURL, err := utils.UploadImageIntoCloudinary(
			p.Cloudinary,
			ctx,
			file,
			"daily-coffee/products",
		)
		if err != nil {
			return fmt.Errorf("failed to upload image: %w", err)
		}
		product.ImageURL = imageURL
	}

	// Update field hanya kalau ada input baru
	if request.Name != "" {
		product.Name = request.Name
		product.Slug = utils.GenerateSlug(request.Name)
		product.SKU = utils.GenerateSKU(request.Name)
	}
	if request.Variant != "" {
		product.Variant = request.Variant
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price != 0 {
		product.Price = request.Price
	}
	if request.Stock != 0 {
		product.Stock = request.Stock
	}
	if request.CategoryID != uuid.Nil {
		product.CategoryID = request.CategoryID
	}

	// check duplicate
	exists, err := p.ProductRepository.ExistsByName(p.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "product name already exist")
	}

	err = p.ProductRepository.Update(p.DB.WithContext(ctx), product)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return utils.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (p *ProductUseCase) Delete(ctx context.Context, productID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	product := &entity.Product{}
	_, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), product, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return utils.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = p.ProductRepository.Delete(p.DB.WithContext(ctx), product)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found, id=%s", productID)
			return utils.ErrNotFound
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (p *ProductUseCase) FindSpecialProduct(ctx context.Context) (*model.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	product, err := p.ProductRepository.FindSpecialProduct(p.DB.WithContext(ctx))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product special is empty")
			return nil, fmt.Errorf("%w: %s", utils.ErrNotFound, err.Error())
		}
		p.Log.Warnf("Failed find special product from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return converter.ProductToResponse(product), nil
}

func (p *ProductUseCase) SetSpecialProduct(ctx context.Context, request *model.UpdateSpecialProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(p.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	products := &entity.Product{}
	product, err := p.ProductRepository.FindById(p.DB.WithContext(ctx), products, request.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Log.Infof("product not found")
			return fmt.Errorf("%w: %s", utils.ErrNotFound, err.Error())
		}
		p.Log.Warnf("Failed find product from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	prevSpecialProduct, err := p.ProductRepository.FindSpecialProduct(p.DB.WithContext(ctx))
	if err == nil && prevSpecialProduct != nil {
		prevSpecialProduct.IsSpecial = false
		if err := p.ProductRepository.Update(p.DB.WithContext(ctx), prevSpecialProduct); err != nil {
			return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	product.IsSpecial = request.IsSpecial
	if err := p.ProductRepository.Update(p.DB.WithContext(ctx), product); err != nil {
		p.Log.Warnf("Failed update special product : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	prevSpecialProduct.IsSpecial = false
	product.IsSpecial = request.IsSpecial

	return nil
}
