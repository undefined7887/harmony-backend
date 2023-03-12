package authdomain

type AuthDTO struct {
	UserID    string `json:"user_id"`
	UserToken string `json:"user_token"`
}

type SignUpRequestBody struct {
	Idtoken  string `json:"idtoken" binding:"jwt"`
	Nickname string `json:"nickname" binding:"nickname"`
}

type SignUpResponse struct {
	AuthDTO
}

type SignInRequestBody struct {
	Idtoken string `json:"idtoken" binding:"jwt"`
}

type SignInResponse struct {
	AuthDTO
}
