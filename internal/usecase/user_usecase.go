package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model/converter"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed count user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("User already exists : %+v", err)
		return nil, fiber.ErrConflict
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       request.ID,
		Password: string(password),
		Name:     request.Name,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}
