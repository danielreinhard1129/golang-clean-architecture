package repository

import "gorm.io/gorm"

type Adapter struct {
	DB               *gorm.DB
	User             UserRepository
	VerificationCode VerificationCodeRepository
}

func NewAdapter(db *gorm.DB) *Adapter {
	return &Adapter{
		DB:               db,
		User:             NewUserRepository(db),
		VerificationCode: NewVerificationCodeRepository(db),
	}
}

func (a *Adapter) WithTx(tx *gorm.DB) *Adapter {
	return &Adapter{
		DB:               tx,
		User:             NewUserRepository(tx),
		VerificationCode: NewVerificationCodeRepository(tx),
	}
}
