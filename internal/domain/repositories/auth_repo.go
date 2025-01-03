package repositories

import (
	"chat-application/users/internal/domain/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

const (
	ErrDuplicatedKey  = "Username or email or phonenumber is already exists."
	ErrRecordNotFound = "User doesn't exists."
	ErrNotUpdated     = "Value is the same"
	ErrDbError        = "Failed to process data"
)

type AuthRepository interface {
	Insert(ctx context.Context, user *models.Mst_users) error
	FindByUsername(ctx context.Context, user *models.Mst_users) (models.Mst_users, error)
	FindByID(ctx context.Context, user *models.Mst_users) (models.Mst_users, error)
	UpdateById(ctx context.Context, id uint, columnName string, updateData interface{}) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Insert(ctx context.Context, user *models.Mst_users) error {
	result := r.db.
		WithContext(ctx).
		Model(&models.Mst_users{}).
		Create(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New(ErrDuplicatedKey)
		}
		return errors.New(ErrDbError)
	}
	return nil
}

func (r *authRepository) FindByUsername(ctx context.Context, user *models.Mst_users) (models.Mst_users, error) {
	var find models.Mst_users
	result := r.db.
		WithContext(ctx).
		Model(&models.Mst_users{}).
		Where("username = ?", user.Username).
		First(&find)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Mst_users{}, errors.New(ErrRecordNotFound)
		}
		return models.Mst_users{}, errors.New(ErrDbError)
	}
	return find, nil
}

func (r *authRepository) FindByID(ctx context.Context, user *models.Mst_users) (models.Mst_users, error) {
	var find models.Mst_users
	result := r.db.
		WithContext(ctx).
		Model(&models.Mst_users{}).
		First(&find, user.ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Mst_users{}, errors.New(ErrRecordNotFound)
		}
		return models.Mst_users{}, errors.New(ErrDbError)
	}
	return find, nil
}

func (r *authRepository) UpdateById(ctx context.Context, id uint, columnName string, updateData interface{}) error {
	var find models.Mst_users
	result := r.db.
		WithContext(ctx).
		Model(&models.Mst_users{}).
		Where("id = ?", id).
		First(&find)

	if find.IsDeleted {
		return errors.New(ErrNotUpdated)
	}

	result.Update(columnName, updateData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New(ErrRecordNotFound)
		}
		return errors.New(ErrDbError)
	}
	return nil
}
