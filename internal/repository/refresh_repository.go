package repository

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RefreshRepository struct {
	Repository[entity.RefreshSession]
	Log *logrus.Logger
}

func NewRefreshRepository(log *logrus.Logger) *RefreshRepository {
	return &RefreshRepository{
		Log: log,
	}
}

func (r *RefreshRepository) GetByJTI(db *gorm.DB, jti string) (*entity.RefreshSession, error) {
	var refresh entity.RefreshSession
	if err := db.Where("jti = ?", jti).Take(&refresh).Error; err != nil {
		return nil, err
	}
	return &refresh, nil
}

// func (r *RefreshRepository) Revoke(db *gorm.DB, jti string) error {
// 	return db.Model(&entity.RefreshSession{}).Where("jti = ?", jti).Update("revoked", true).Error
// }

// func (r *RefreshRepository) RevokeUserAll(db *gorm.DB, userID string) error {
// 	return db.Model(&entity.RefreshSession{}).Where("user_id = ?", userID).Update("revoked", true).Error
// }

// func (r *RefreshRepository) CleanupExpired(db *gorm.DB, now time.Time) error {
// 	return db.Where("expires_at < ?", now).Delete(&entity.RefreshSession{}).Error
// }
