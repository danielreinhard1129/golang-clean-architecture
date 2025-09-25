package response

import "github.com/danielreinhard1129/fiber-clean-arch/internal/entities"

type AuthLoginResponse struct {
	entities.User
	Token string `json:"token"`
}

type AuthRegisterResponse struct {
	Message string `json:"message"`
}

type AuthVerifyAccountResponse struct {
	Message string `json:"message"`
}
