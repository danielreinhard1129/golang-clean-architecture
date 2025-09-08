package repository

import (
	"errors"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(search string, offset, limit int) ([]entities.User, int64)
	FindById(id int) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	Create(user entities.User) (entities.User, error)
}

type userRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) FindAll(search string, offset, limit int) ([]entities.User, int64) {
	var users []entities.User
	var total int64

	query := r.DB.Model(&entities.User{})

	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return []entities.User{}, 0
	}

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return []entities.User{}, 0
	}

	return users, total
}

func (r *userRepositoryImpl) FindById(id int) (entities.User, error) {
	var user entities.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, exception.NotFoundError{Message: "user not found"}
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, exception.NotFoundError{Message: "email not found"}
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) Create(user entities.User) (entities.User, error) {
	err := r.DB.Create(&user).Error
	return user, err
}
