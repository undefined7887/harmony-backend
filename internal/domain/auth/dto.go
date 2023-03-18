package authdomain

import userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"

type AuthDTO struct {
	User      userdomain.UserDTO `json:"user"`
	UserToken string             `json:"user_token"`
}

type SignUpRequestBody struct {
	Nonce    string `json:"nonce" binding:"hexadecimal"`
	Idtoken  string `json:"idtoken" binding:"jwt"`
	Nickname string `json:"nickname" binding:"nickname"`
}

type SignUpResponse struct {
	AuthDTO
}

type SignInRequestBody struct {
	Nonce   string `json:"nonce" binding:"hexadecimal"`
	Idtoken string `json:"idtoken" binding:"jwt"`
}

type SignInResponse struct {
	AuthDTO
}
