package converter

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
