package repository

import (
	"context"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"gorm.io/gorm"
)

type VerificationCodeRepository interface {
	Create(ctx context.Context, vc entities.VerificationCode) (entities.VerificationCode, error)
	FindLatestByUserAndPurpose(ctx context.Context, userID uint, purpose string) (entities.VerificationCode, error)
	DeleteAllByUserAndPurpose(ctx context.Context, userID uint, purpose string) error
}

type verificationCodeRepositoryImpl struct {
	db *gorm.DB
}

func NewVerificationCodeRepository(db *gorm.DB) VerificationCodeRepository {
	return &verificationCodeRepositoryImpl{db}
}

func (r *verificationCodeRepositoryImpl) Create(ctx context.Context, vc entities.VerificationCode) (entities.VerificationCode, error) {
	err := r.db.WithContext(ctx).Create(&vc).Error
	return vc, err
}

func (r *verificationCodeRepositoryImpl) FindLatestByUserAndPurpose(ctx context.Context, userID uint, purpose string) (entities.VerificationCode, error) {
	var vc entities.VerificationCode
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose = ?", userID, purpose).
		Order("created_at DESC").
		First(&vc).Error
	if err != nil {
		return entities.VerificationCode{}, err
	}
	return vc, nil
}

func (r *verificationCodeRepositoryImpl) DeleteAllByUserAndPurpose(ctx context.Context, userID uint, purpose string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND purpose = ?", userID, purpose).
		Delete(&entities.VerificationCode{}).Error
}
