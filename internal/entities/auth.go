package entities

type AuthResponse struct {
	User
	Token string `json:"token"`
}
