package repository

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Repository[entity.Customer]
	Log *logrus.Logger
}

func NewCustomerRepository(log *logrus.Logger) *CustomerRepository {
	return &CustomerRepository{
		Log: log,
	}
}

func (r *CustomerRepository) CountByUserName(db *gorm.DB, userName string) (int64, error) {
	var count int64
	err := db.Model(&entity.Customer{}).Where("user_name = ?", userName).Count(&count).Error
	return count, err
}

func (r *CustomerRepository) ExistsByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	err := db.Model(&entity.Customer{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
