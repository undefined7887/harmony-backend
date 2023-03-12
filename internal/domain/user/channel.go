package userdomain

import "fmt"

const (
	ChannelNamespace = "user"
)

func ChannelUser(userID string) string {
	return fmt.Sprintf("%s:%s", ChannelNamespace, userID)
}
