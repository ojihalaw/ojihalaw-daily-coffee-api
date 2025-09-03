package converter

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
)

func CustomerToResponse(user *entity.Customer) *model.CustomerResponse {
	return &model.CustomerResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
