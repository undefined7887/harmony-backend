package userdomain

type UserDTO struct {
	ID string `json:"id"`

	Photo    string `json:"photo"`
	Nickname string `json:"nickname"`
}
