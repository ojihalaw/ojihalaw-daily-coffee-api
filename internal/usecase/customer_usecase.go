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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validator          *utils.Validator
	CustomerRepository *repository.CustomerRepository
}

func NewCustomerUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	customerRepository *repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		DB:                 db,
		Log:                logger,
		Validator:          validator,
		CustomerRepository: customerRepository,
	}
}

func (c *CustomerUseCase) Create(ctx context.Context, request *model.RegisterCustomerRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

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
	total, err := c.CustomerRepository.CountByUserName(tx, request.UserName)
	if err != nil {
		c.Log.Warnf("Failed count user from database: %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}
	if total > 0 {
		c.Log.Warnf("User already exists : %+v", request.UserName)
		return fmt.Errorf("%w: %s", utils.ErrConflict, "user with username already exists")
	}

	exists, err := c.CustomerRepository.ExistsByEmail(tx, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "email already registered")
	}

	// hash password
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	user := &entity.Customer{
		Password:    string(password),
		Name:        request.Name,
		UserName:    request.UserName,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
	}

	if err := c.CustomerRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *CustomerUseCase) FindAll(ctx context.Context, pagination *utils.PaginationRequest) ([]model.CustomerResponse, *utils.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var customers []entity.Customer

	total, err := c.CustomerRepository.FindAll(c.DB.WithContext(ctx), &customers, pagination)
	if err != nil {
		c.Log.Warnf("Failed find all customer from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	responses := make([]model.CustomerResponse, len(customers))
	for i, customer := range customers {
		responses[i] = *converter.CustomerToResponse(&customer)
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

func (c *CustomerUseCase) FindByID(ctx context.Context, customerID string) (*model.CustomerResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var customers *entity.Customer

	customer, err := c.CustomerRepository.FindById(c.DB.WithContext(ctx), customers, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("customer not found, id=%s", customerID)
			return nil, utils.ErrNotFound
		}
		c.Log.Warnf("Failed find customer from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Update(ctx context.Context, customerID string, request *model.UpdateCustomerRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	customer := &entity.Customer{}
	_, err := c.CustomerRepository.FindById(c.DB.WithContext(ctx), customer, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("customer not found, id=%s", customerID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find customer from database : %+v", err)
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

	total, err := c.CustomerRepository.CountByUserName(c.DB.WithContext(ctx), request.UserName)
	if err != nil {
		c.Log.Warnf("Failed count user from database: %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}
	if total > 0 {
		c.Log.Warnf("User already exists : %+v", request.UserName)
		return fmt.Errorf("%w: %s", utils.ErrConflict, "user with username already exists")
	}

	exists, err := c.CustomerRepository.ExistsByEmail(c.DB.WithContext(ctx), request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "email already registered")
	}

	customer.Name = request.Name

	err = c.CustomerRepository.Update(c.DB.WithContext(ctx), customer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("customer not found, id=%s", customerID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find customer from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *CustomerUseCase) Delete(ctx context.Context, customerID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	customer := &entity.Customer{}
	_, err := c.CustomerRepository.FindById(c.DB.WithContext(ctx), customer, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("customer not found, id=%s", customerID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find customer from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.CustomerRepository.Delete(c.DB.WithContext(ctx), customer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("customer not found, id=%s", customerID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find customer from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
