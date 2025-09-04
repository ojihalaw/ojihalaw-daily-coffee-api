package repository

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	Repository[entity.Category]
	Log *logrus.Logger
}

func NewCategoryRepository(log *logrus.Logger) *CategoryRepository {
	return &CategoryRepository{
		Log: log,
	}
}

func (r *CategoryRepository) ExistsByName(db *gorm.DB, name string) (bool, error) {
	var count int64
	err := db.Model(&entity.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
