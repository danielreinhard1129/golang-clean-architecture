package entities

type AuthLogin struct {
	User
	Token string `json:"token"`
}
