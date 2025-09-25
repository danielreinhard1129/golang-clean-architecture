package repository

import "gorm.io/gorm"

type Adapter struct {
	DB    *gorm.DB
	User  UserRepository
	Token TokenRepository
}

func NewAdapter(db *gorm.DB) *Adapter {
	return &Adapter{
		DB:    db,
		User:  NewUserRepository(db),
		Token: NewTokenRepository(db),
	}
}

func (a *Adapter) WithTx(tx *gorm.DB) *Adapter {
	return &Adapter{
		DB:    tx,
		User:  NewUserRepository(tx),
		Token: NewTokenRepository(tx),
	}
}
