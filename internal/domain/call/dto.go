package calldomain

import (
	"encoding/json"
)

type CallDTO struct {
	ID string `json:"_id"`

	UserID string `json:"user_id"`
	PeerID string `json:"peer_id"`

	Status string `json:"status"`

	// WebRTC data
	WebRTC CallWebRtcDTO `json:"web_rtc"`
}

func MapCallDTO(call Call) CallDTO {
	return CallDTO{
		ID:     call.ID,
		UserID: call.UserID,
		PeerID: call.PeerID,
		Status: call.Status,
		WebRTC: MapCallWebRtcDTO(call.WebRTC),
	}
}

type CallWebRtcDTO struct {
	Offer  json.RawMessage `json:"offer"`
	Answer json.RawMessage `json:"answer"`
}

func MapCallWebRtcDTO(callWebRtc CallWebRTC) CallWebRtcDTO {
	return CallWebRtcDTO{
		Offer:  callWebRtc.Offer,
		Answer: callWebRtc.Answer,
	}
}

type PeerParams struct {
	PeerID string `uri:"peer_id" binding:"id"`
}

// --

type CreateCallRequestBody struct {
	WebRtcOffer json.RawMessage `json:"web_rtc_offer"`
}

type CreateCallResponse struct {
	CallID string `json:"call_id"`
}

type NewCallNotification struct {
	CallDTO
}

// --

type UpdateCallRequestBody struct {
	Accept       bool            `json:"accept"`
	WebRtcAnswer json.RawMessage `json:"web_rtc_answer"`
}

type UpdateCallNotification struct {
	CallDTO
}

// --

type ProxyCallDataRequestBody struct {
	WebRtcCandidate json.RawMessage `json:"web_rtc_candidate"`
}

type CallDataNotification struct {
	ID string `json:"id"`

	WebRtcCandidate json.RawMessage `json:"web_rtc_candidate"`
}
