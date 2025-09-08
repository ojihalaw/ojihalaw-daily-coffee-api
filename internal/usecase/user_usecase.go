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

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validator      *utils.Validator
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validator:      validator,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) error {
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

	total, err := c.UserRepository.CountByUserName(tx, request.UserName)
	if err != nil {
		c.Log.Warnf("Failed count user from database: %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}
	if total > 0 {
		c.Log.Warnf("User already exists : %+v", request.UserName)
		return fmt.Errorf("%w: %s", utils.ErrConflict, "user with username already exists")
	}

	exists, err := c.UserRepository.ExistsByEmail(tx, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "email already registered")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	user := &entity.User{
		Password:    string(password),
		Name:        request.Name,
		UserName:    request.UserName,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
		Role:        entity.Role(request.Role),
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *UserUseCase) FindAll(ctx context.Context, pagination *utils.PaginationRequest) ([]model.UserResponse, *utils.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var users []entity.User

	total, err := c.UserRepository.FindAll(c.DB.WithContext(ctx), &users, pagination)
	if err != nil {
		c.Log.Warnf("Failed find all users from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	responses := make([]model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *converter.UserToResponse(&user)
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

func (c *UserUseCase) FindByID(ctx context.Context, userID string) (*model.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var users *entity.User

	user, err := c.UserRepository.FindById(c.DB.WithContext(ctx), users, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("user not found, id=%s", userID)
			return nil, utils.ErrNotFound
		}
		c.Log.Warnf("Failed find user from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Update(ctx context.Context, userID string, request *model.UpdateUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &entity.User{}
	_, err := c.UserRepository.FindById(c.DB.WithContext(ctx), user, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Category not found, id=%s", userID)
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

	user.Name = request.Name
	user.Role = entity.Role(request.Role)

	err = c.UserRepository.Update(c.DB.WithContext(ctx), user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("user not found, id=%s", userID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find user from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *UserUseCase) Delete(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &entity.User{}
	_, err := c.UserRepository.FindById(c.DB.WithContext(ctx), user, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("user not found, id=%s", userID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find user from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.UserRepository.Delete(c.DB.WithContext(ctx), user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("user not found, id=%s", userID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find user from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
