package chatdomain

import "fmt"

func ChannelNewMessage(userID string) string {
	return fmt.Sprintf("chat/message/new#%s", userID)
}

func ChannelUpdatedMessage(userID string) string {
	return fmt.Sprintf("chat/message/updated#%s", userID)
}

func ChannelReadMessage(userID string) string {
	return fmt.Sprintf("chat/message/read#%s", userID)
}

func ChannelTyping(userID string) string {
	return fmt.Sprintf("chat/typing#%s", userID)
}
