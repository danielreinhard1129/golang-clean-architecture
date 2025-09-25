package repository

import (
	"context"
	"errors"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context, search, orderBy, sort string, offset, limit int) ([]entities.User, int64)
	FindById(ctx context.Context, id int) (entities.User, error)
	FindByEmail(ctx context.Context, email string) (entities.User, error)
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Update(ctx context.Context, id int, updatedUser entities.User) (entities.User, error)
	Delete(ctx context.Context, id int) error
}

type userRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, search, orderBy, sort string, offset, limit int) ([]entities.User, int64) {
	var users []entities.User
	var total int64

	query := r.DB.WithContext(ctx).Model(&entities.User{})

	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return []entities.User{}, 0
	}

	if orderBy != "" {
		if sort != "asc" && sort != "desc" {
			sort = "desc"
		}
		query = query.Order(orderBy + " " + sort)
	}

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return []entities.User{}, 0
	}

	return users, total
}

func (r *userRepositoryImpl) FindById(ctx context.Context, id int) (entities.User, error) {
	var user entities.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, exception.NotFoundError{Message: "user not found"}
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, exception.NotFoundError{Message: "email not found"}
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, user entities.User) (entities.User, error) {
	err := r.DB.WithContext(ctx).Create(&user).Error
	return user, err
}

func (r *userRepositoryImpl) Update(ctx context.Context, id int, updatedUser entities.User) (entities.User, error) {
	var user entities.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, exception.NotFoundError{Message: "user not found"}
		}
		return entities.User{}, err
	}

	if err := r.DB.WithContext(ctx).Model(&user).Updates(updatedUser).Error; err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int) error {
	var user entities.User

	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NotFoundError{Message: "user not found"}
		}
		return err
	}

	if err := r.DB.WithContext(ctx).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
