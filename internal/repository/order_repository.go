package repository

import (
	"context"

	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Repository[entity.Order]
	Log *logrus.Logger
}

func NewOrderRepository(log *logrus.Logger) *OrderRepository {
	return &OrderRepository{
		Log: log,
	}
}

func (r *OrderRepository) GetTodayOrderCount(ctx context.Context, tx *gorm.DB) (int, error) {
	var count int64
	err := tx.Model(&entity.Order{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&count).Error
	return int(count), err
}
