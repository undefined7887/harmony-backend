package calldomain

import (
	"fmt"
)

const (
	ChannelNamespace = "call"
)

func ChannelCallNew(userID string) string {
	return fmt.Sprintf("%s:new#%s", ChannelNamespace, userID)
}

func ChannelCallUpdates(userID string) string {
	return fmt.Sprintf("%s:updates#%s", ChannelNamespace, userID)
}

func ChannelCallData(userID string) string {
	return fmt.Sprintf("%s:data#%s", ChannelNamespace, userID)
}
