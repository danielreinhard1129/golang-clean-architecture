package usecase

import (
	"errors"
	"time"

	"github.com/danielreinhard1129/fiber-clean-arch/configs"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/repository"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(reqBody *request.AuthLoginRequest) (entities.AuthResponse, error)
}

type authUsecaseImpl struct {
	userRepo repository.UserRepository
	config   configs.Config
}

func NewAuthUsecase(userRepo *repository.UserRepository, config configs.Config) AuthUsecase {
	return &authUsecaseImpl{userRepo: *userRepo, config: config}
}

func (u *authUsecaseImpl) Login(reqBody *request.AuthLoginRequest) (entities.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(reqBody.Email)
	if err != nil {
		return entities.AuthResponse{}, exception.NotFoundError{Message: "invalid credentials"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return entities.AuthResponse{}, exception.NotFoundError{Message: "invalid credentials"}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // exp 3 days
	})

	tokenString, err := token.SignedString([]byte(u.config.Get("JWT_SECRET")))
	if err != nil {
		return entities.AuthResponse{}, errors.New("failed to generate token")
	}

	return entities.AuthResponse{
		User:  user,
		Token: tokenString,
	}, nil
}
