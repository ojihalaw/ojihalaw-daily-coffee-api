package repository

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Repository[entity.Product]
	Log *logrus.Logger
}

func NewProductRepository(log *logrus.Logger) *ProductRepository {
	return &ProductRepository{
		Log: log,
	}
}

func (r *ProductRepository) ExistsByName(db *gorm.DB, name string) (bool, error) {
	var count int64
	err := db.Model(&entity.Product{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *ProductRepository) FindSpecialProduct(db *gorm.DB) (*entity.Product, error) {
	var p entity.Product
	if err := db.Where("is_special = ?", true).Take(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
