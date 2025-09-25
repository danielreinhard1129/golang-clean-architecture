package repository

import (
	"context"
	"errors"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"gorm.io/gorm"
)

type TokenRepository interface {
	FindByToken(ctx context.Context, token string) (entities.Token, error)
	Create(ctx context.Context, token entities.Token) (entities.Token, error)
}

type tokenRepositoryImpl struct {
	*gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepositoryImpl{DB: db}
}

func (r *tokenRepositoryImpl) FindByToken(ctx context.Context, token string) (entities.Token, error) {
	var result entities.Token
	if err := r.DB.WithContext(ctx).
		Where("token = ?", token).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Token{}, exception.NotFoundError{Message: "token not found"}
		}
		return entities.Token{}, err
	}
	return result, nil
}

func (r *tokenRepositoryImpl) Create(ctx context.Context, token entities.Token) (entities.Token, error) {
	err := r.DB.WithContext(ctx).Create(&token).Error
	return token, err
}
