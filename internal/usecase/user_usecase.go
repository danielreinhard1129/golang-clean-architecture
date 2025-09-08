package usecase

import (
	"errors"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/repository"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	FindAll(search, orderBy, sort string, page, limit int) ([]entities.User, int64)
	FindById(id int) (entities.User, error)
	Create(reqBody *request.UserCreateRequest) (entities.User, error)
	Update(id int, reqBody *request.UserUpdateRequest) (entities.User, error)
	Delete(id int) error
}

type userUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repository *repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo: *repository}
}

func (u *userUsecaseImpl) FindAll(search, orderBy, sort string, page, limit int) ([]entities.User, int64) {
	offset := (page - 1) * limit
	users, total := u.repo.FindAll(search, orderBy, sort, offset, limit)
	return users, total
}

func (u *userUsecaseImpl) FindById(id int) (entities.User, error) {
	return u.repo.FindById(id)
}

func (u *userUsecaseImpl) Create(reqBody *request.UserCreateRequest) (entities.User, error) {

	_, err := u.repo.FindByEmail(reqBody.Email)

	if err == nil {
		return entities.User{}, exception.ConflictError{Message: "email already exists"}
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, errors.New("error hashing password")
	}
	hashedPassword := string(hashedPasswordBytes)

	user := entities.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: hashedPassword,
	}

	return u.repo.Create(user)
}

func (u *userUsecaseImpl) Update(id int, reqBody *request.UserUpdateRequest) (entities.User, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return entities.User{}, exception.NotFoundError{Message: "user not found"}
	}

	if reqBody.Name != "" {
		user.Name = reqBody.Name
	}

	if reqBody.Email != "" && reqBody.Email != user.Email {
		_, err := u.repo.FindByEmail(reqBody.Email)
		if err == nil {
			return entities.User{}, exception.ConflictError{Message: "email already exists"}
		}
		user.Email = reqBody.Email
	}

	if reqBody.Password != "" {
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
		if err != nil {
			return entities.User{}, errors.New("error hashing password")
		}
		user.Password = string(hashedPasswordBytes)
	}

	updatedUser, err := u.repo.Update(id, user)
	if err != nil {
		return entities.User{}, err
	}

	return updatedUser, nil
}

func (u *userUsecaseImpl) Delete(id int) error {
	_, err := u.repo.FindById(id)
	if err != nil {
		return exception.NotFoundError{Message: "user not found"}
	}

	if err := u.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
