package userdomain

type UserProfileDTO struct {
	ID       string `json:"id"`
	Photo    string `json:"photo"`
	Nickname string `json:"nickname"`
}
