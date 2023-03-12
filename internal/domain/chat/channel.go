package chatdomain

import "fmt"

const (
	ChannelMessageCreated = "message_created"
	ChannelMessageUpdated = "message_updated"

	ChannelRead   = "read"
	ChannelTyping = "typing"
)

func Channel(channel, userID string) string {
	return fmt.Sprintf("chat:%s#%s", channel, userID)
}
