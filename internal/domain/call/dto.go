package calldomain

import (
	"encoding/json"
)

type CallDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	PeerID string `json:"peer_id"`
	Status string `json:"status"`
}

func MapCallDTO(call Call) CallDTO {
	return CallDTO{
		ID:     call.ID,
		UserID: call.UserID,
		PeerID: call.PeerID,
		Status: call.Status,
	}
}

type PeerParams struct {
	PeerID string `uri:"peer_id" binding:"id"`
}

// --

type CreateCallResponse struct {
	CallID string `json:"call_id"`
}

type NewCallNotification struct {
	CallDTO
}

// --

type UpdateCallRequestBody struct {
	Status string `json:"status" binding:"oneof=accepted declined finished"`
}

type UpdateCallNotification struct {
	CallDTO
}

// --

type ProxyCallDataRequestBody struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type CallDataNotification struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}
