package userdomain

import "github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"

type UserDTO struct {
	ID       string `json:"id"`
	Status   string `json:"status,omitempty"`
	Photo    string `json:"photo"`
	Nickname string `json:"nickname"`
}

func MapUserDTO(user User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		Status:   user.Status,
		Photo:    user.Photo,
		Nickname: user.Nickname,
	}
}

type UpdateUserNotification struct {
	UserDTO
}

// ---

type GetUserRequestParams struct {
	ID string `uri:"id" binding:"id|eq=self"`
}

type GetUserResponse struct {
	UserDTO
}

// ---

type GetUserByNicknameRequestQuery struct {
	Nickname string `form:"nickname" binding:"nickname-extended"`
}

type GetUserByNicknameResponse struct {
	UserDTO
}

// ---

type UpdateUserStatusRequestBody struct {
	Status string `json:"status" binding:"oneof=online away silence"`
}

// ---

var CentrifugoUnauthorizedResponse = &centrifugo.Response[any]{
	Error: &centrifugo.ResponseError{
		Code:    4500, // Code 4500 doesn't allow client to reconnect
		Message: "Unauthorized",
	},
}

// ---

type CentrifugoConnectResponse struct {
	User     string `json:"user"`
	ExpireAt int64  `json:"expire_at"`
}

// ---

type CentrifugoRefreshResponse struct {
	ExpireAt int64 `json:"expire_at"`
}
