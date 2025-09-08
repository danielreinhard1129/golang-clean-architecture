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
	FindAll(search string, page, limit int) ([]entities.User, int64)
	FindById(id int) (entities.User, error)
	Create(reqBody *request.UserCreateRequest) (entities.User, error)
}

type userUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repository *repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo: *repository}
}

func (u *userUsecaseImpl) FindAll(search string, page, limit int) ([]entities.User, int64) {
	offset := (page - 1) * limit
	users, total := u.repo.FindAll(search, offset, limit)
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
