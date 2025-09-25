package usecase

import (
	"context"
	"errors"

	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/repository"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/mail"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	FindAll(ctx context.Context, search, orderBy, sort string, page, limit int) ([]entities.User, int64)
	FindById(ctx context.Context, id int) (entities.User, error)
	Create(ctx context.Context, reqBody *request.UserCreateRequest) (entities.User, error)
	Update(ctx context.Context, id int, reqBody *request.UserUpdateRequest) (entities.User, error)
	Delete(ctx context.Context, id int) error
}

type userUsecaseImpl struct {
	adapter     repository.Adapter
	mailService mail.Service
}

func NewUserUsecase(adapter *repository.Adapter, mailService *mail.Service) UserUsecase {
	return &userUsecaseImpl{adapter: *adapter, mailService: *mailService}
}

func (u *userUsecaseImpl) FindAll(ctx context.Context, search, orderBy, sort string, page, limit int) ([]entities.User, int64) {
	offset := (page - 1) * limit
	users, total := u.adapter.User.FindAll(ctx, search, orderBy, sort, offset, limit)
	return users, total
}

func (u *userUsecaseImpl) FindById(ctx context.Context, id int) (entities.User, error) {
	return u.adapter.User.FindById(ctx, id)
}

func (u *userUsecaseImpl) Create(ctx context.Context, reqBody *request.UserCreateRequest) (entities.User, error) {

	_, err := u.adapter.User.FindByEmail(ctx, reqBody.Email)

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

	go func() {
		err := u.mailService.SendMail(
			user.Email,
			"Welcome to MyApp",
			"welcome.html",
			map[string]any{"Name": user.Name},
		)
		if err != nil {
			println("failed to send email:", err.Error())
		}
	}()

	return u.adapter.User.Create(ctx, user)
}

func (u *userUsecaseImpl) Update(ctx context.Context, id int, reqBody *request.UserUpdateRequest) (entities.User, error) {
	user, err := u.adapter.User.FindById(ctx, id)
	if err != nil {
		return entities.User{}, exception.NotFoundError{Message: "user not found"}
	}

	if reqBody.Name != "" {
		user.Name = reqBody.Name
	}

	if reqBody.Email != "" && reqBody.Email != user.Email {
		_, err := u.adapter.User.FindByEmail(ctx, reqBody.Email)
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

	updatedUser, err := u.adapter.User.Update(ctx, id, user)
	if err != nil {
		return entities.User{}, err
	}

	return updatedUser, nil
}

func (u *userUsecaseImpl) Delete(ctx context.Context, id int) error {
	_, err := u.adapter.User.FindById(ctx, id)
	if err != nil {
		return exception.NotFoundError{Message: "user not found"}
	}

	if err := u.adapter.User.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
