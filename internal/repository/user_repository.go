package repository

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) CountByUserName(db *gorm.DB, userName string) (int64, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("user_name = ?", userName).Count(&count).Error
	return count, err
}

func (r *UserRepository) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	var u entity.User
	if err := db.Where("email = ?", email).Take(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
