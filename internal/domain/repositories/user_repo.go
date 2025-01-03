package repositories

import (
	"chat-application/users/internal/domain/dtos"
	"chat-application/users/internal/domain/models"
	"chat-application/users/pkg/database"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(ctx context.Context, page int, limit int) ([]dtos.GetAllUserResponsePayload, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAll(ctx context.Context, page int, limit int) ([]dtos.GetAllUserResponsePayload, int64, error) {
	var users []dtos.GetAllUserResponsePayload
	var count int64
	result := r.db.
		WithContext(ctx).
		Model(&models.Mst_users{}).
		Select("id", "username", "phonenumber", "email", "is_online").
		Where("is_deleted", false).
		Count(&count).
		Scopes(database.Paginate(page, limit)).
		Find(&users)

	if result.Error != nil {
		return []dtos.GetAllUserResponsePayload{}, 0, errors.New(ErrDbError)
	}

	return users, count, nil
}
