package migration

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, log *logrus.Logger) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Customer{},
		&entity.Category{},
		&entity.Product{},
		&entity.RefreshSession{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.PaymentLog{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Info("Migration success ✅")
}
