package chatdomain

import "fmt"

const (
	ChannelNamespace = "chat"
)

func ChannelMessageNew(userID string) string {
	return fmt.Sprintf("%s:message/new#%s", ChannelNamespace, userID)
}

func ChannelMessageUpdates(userID string) string {
	return fmt.Sprintf("%s:message/updates#%s", ChannelNamespace, userID)
}

func ChannelReadUpdates(userID string) string {
	return fmt.Sprintf("%s:read/updates#%s", ChannelNamespace, userID)
}

func ChannelTypingUpdates(userID string) string {
	return fmt.Sprintf("%s:typing/updates#%s", ChannelNamespace, userID)
}
