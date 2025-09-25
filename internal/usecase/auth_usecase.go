package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/danielreinhard1129/fiber-clean-arch/configs"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/delivery/http/request"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/repository"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/mail"
	"github.com/danielreinhard1129/fiber-clean-arch/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Login(ctx context.Context, reqBody *request.AuthLoginRequest) (entities.AuthLogin, error)
	Register(ctx context.Context, reqBody *request.AuthRegisterRequest) error
	VerifyAccount(ctx context.Context, reqBody *request.AuthVerifyAccountRequest) error
	ResendOTP(ctx context.Context, reqBody *request.AuthResendOTPRequest) error
}

type authUsecaseImpl struct {
	adapter     repository.Adapter
	config      configs.Config
	mailService mail.Service
}

func NewAuthUsecase(adapter *repository.Adapter, mailService *mail.Service, config configs.Config) AuthUsecase {
	return &authUsecaseImpl{adapter: *adapter, mailService: *mailService, config: config}
}

func (u *authUsecaseImpl) Login(ctx context.Context, reqBody *request.AuthLoginRequest) (entities.AuthLogin, error) {
	user, err := u.adapter.User.FindByEmail(ctx, reqBody.Email)
	if err != nil {
		return entities.AuthLogin{}, exception.NotFoundError{Message: "invalid credentials"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return entities.AuthLogin{}, exception.NotFoundError{Message: "invalid credentials"}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // exp 3 days
	})

	tokenString, err := token.SignedString([]byte(u.config.Get("JWT_SECRET")))
	if err != nil {
		return entities.AuthLogin{}, errors.New("failed to generate token")
	}

	return entities.AuthLogin{
		User:  user,
		Token: tokenString,
	}, nil
}

func (u *authUsecaseImpl) Register(ctx context.Context, reqBody *request.AuthRegisterRequest) error {
	return u.adapter.DB.Transaction(func(tx *gorm.DB) error {
		txAdapter := u.adapter.WithTx(tx)

		_, err := txAdapter.User.FindByEmail(ctx, reqBody.Email)
		if err == nil {
			return exception.ConflictError{Message: "email already exists"}
		}

		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("error hashing password")
		}
		hashedPassword := string(hashedPasswordBytes)

		user := entities.User{
			Name:     reqBody.Name,
			Email:    reqBody.Email,
			Password: hashedPassword,
		}

		user, err = txAdapter.User.Create(ctx, user)
		if err != nil {
			return err
		}

		otp, err := utils.GenerateOTP(6)
		if err != nil {
			return errors.New("failed to generate otp")
		}

		const otp_expired_time = 15

		verificationCode := entities.VerificationCode{
			UserID:    user.ID,
			Code:      otp,
			ExpiredAt: time.Now().Add(otp_expired_time * time.Minute),
			Purpose:   "EMAIL_VERIFICATION",
		}

		vc, err := txAdapter.VerificationCode.Create(ctx, verificationCode)
		if err != nil {
			return err
		}

		go func() {
			err := u.mailService.SendMail(
				user.Email,
				"Verify Your Account",
				"verify-account.html",
				map[string]any{"ExpiredMinutes": otp_expired_time, "OTP": vc.Code},
			)
			if err != nil {
				println("failed to send email:", err.Error())
			}
		}()

		return nil
	})
}

func (u *authUsecaseImpl) VerifyAccount(ctx context.Context, reqBody *request.AuthVerifyAccountRequest) error {
	return u.adapter.DB.Transaction(func(tx *gorm.DB) error {
		txAdapter := u.adapter.WithTx(tx)

		user, err := txAdapter.User.FindByEmail(ctx, reqBody.Email)
		if err != nil {
			return err
		}

		vc, err := txAdapter.VerificationCode.FindLatestByUserAndPurpose(ctx, user.ID, "EMAIL_VERIFICATION")
		if err != nil {
			return exception.NotFoundError{Message: "verification code not found"}
		}

		if time.Now().After(vc.ExpiredAt) {
			return exception.BadRequestError{Message: "OTP expired"}
		}

		if vc.Code != reqBody.Otp {
			return exception.BadRequestError{Message: "Invalid OTP"}
		}

		now := time.Now()
		_, err = txAdapter.User.Update(ctx, int(user.ID), entities.User{VerifiedAt: &now})
		if err != nil {
			return err
		}

		err = txAdapter.VerificationCode.DeleteAllByUserAndPurpose(ctx, user.ID, "EMAIL_VERIFICATION")
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *authUsecaseImpl) ResendOTP(ctx context.Context, reqBody *request.AuthResendOTPRequest) error {
	return u.adapter.DB.Transaction(func(tx *gorm.DB) error {
		txAdapter := u.adapter.WithTx(tx)

		user, err := txAdapter.User.FindByEmail(ctx, reqBody.Email)
		if err != nil {
			return exception.NotFoundError{Message: "user not found"}
		}

		if user.VerifiedAt != nil {
			return exception.BadRequestError{Message: "account already verified"}
		}

		err = txAdapter.VerificationCode.DeleteAllByUserAndPurpose(ctx, user.ID, "EMAIL_VERIFICATION")
		if err != nil {
			return err
		}

		otp, err := utils.GenerateOTP(6)
		if err != nil {
			return errors.New("failed to generate otp")
		}

		const otp_expired_time = 15
		verificationCode := entities.VerificationCode{
			UserID:    user.ID,
			Code:      otp,
			ExpiredAt: time.Now().Add(otp_expired_time * time.Minute),
			Purpose:   "EMAIL_VERIFICATION",
		}

		vc, err := txAdapter.VerificationCode.Create(ctx, verificationCode)
		if err != nil {
			return err
		}

		go func() {
			err := u.mailService.SendMail(
				user.Email,
				"Verify Your Account",
				"verify-account.html",
				map[string]any{"ExpiredMinutes": otp_expired_time, "OTP": vc.Code},
			)
			if err != nil {
				println("failed to send email:", err.Error())
			}
		}()

		return nil
	})
}
