package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB        *gorm.DB
	Log       *logrus.Logger
	Validator *utils.Validator
	JWT       *utils.JWTMaker
	UserRepo  *repository.UserRepository
	SessRepo  *repository.RefreshRepository
}

func NewAuthUseCase(
	db *gorm.DB, logger *logrus.Logger, validator *utils.Validator, jwt *utils.JWTMaker,
	userRepository *repository.UserRepository, sessRepository *repository.RefreshRepository) *AuthUseCase {
	return &AuthUseCase{
		DB:        db,
		Log:       logger,
		Validator: validator,
		UserRepo:  userRepository,
		SessRepo:  sessRepository,
		JWT:       jwt,
	}
}

func randJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (c *AuthUseCase) Login(ctx context.Context, request *model.LoginRequest, ip, ua string) (*model.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(c.Validator.Translator))
			}
			return nil, fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return nil, fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	user, err := c.UserRepo.FindByEmail(c.DB.WithContext(ctx), request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("user not found, id=%s", request.Email)
			return nil, utils.ErrInvalidEmail
		}
		c.Log.Warnf("Failed find user from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		return nil, utils.ErrInvalidPassword
	}

	jti, err := randJTI()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	// create refresh session
	sess := &entity.RefreshSession{
		UserID: user.ID, JTI: jti,
		ExpiresAt: now.Add(c.JWT.RefreshTTL),
		IP:        ip,
		UserAgent: ua,
	}

	if err := c.SessRepo.Create(c.DB.WithContext(ctx), sess); err != nil {
		return nil, err
	}

	access, accessExp, err := c.JWT.NewAccessToken(user.ID.String(), string(user.Role), now)
	if err != nil {
		return nil, err
	}

	refresh, refreshExp, err := c.JWT.NewRefreshToken(user.ID.String(), jti, now)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:      access,
		RefreshToken:     refresh,
		AccessExpiresIn:  int64(time.Until(accessExp).Seconds()),
		RefreshExpiresIn: int64(time.Until(refreshExp).Seconds()),
	}, nil
}
