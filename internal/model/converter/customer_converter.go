package converter

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
)

func CustomerToResponse(user *entity.Customer) *model.CustomerResponse {
	return &model.CustomerResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		UserName:    user.UserName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
	}
}
