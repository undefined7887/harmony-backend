package chattransport

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	chatservice "github.com/undefined7887/harmony-backend/internal/service/chat"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/transport"
	authtransport "github.com/undefined7887/harmony-backend/internal/transport/auth"
	"net/http"
)

type HttpEndpoint struct {
	service    *chatservice.Service
	jwtService *jwtservice.Service
}

func NewHttpEndpoint(service *chatservice.Service, jwtService *jwtservice.Service) transport.HttpEndpoint {
	return &HttpEndpoint{
		service:    service,
		jwtService: jwtService,
	}
}

func (e *HttpEndpoint) Register(group *gin.RouterGroup) {
	chatGroup := group.
		Group("/chat").
		Use(authtransport.NewHttpAuthMiddleware(e.jwtService))
	{
		chatGroup.POST("/:peer_type/:peer_id", e.createMessage)

		chatGroup.GET("", e.listMessages) // replace function by listChats
		chatGroup.GET("/:peer_type", e.listMessages)
		chatGroup.GET("/:peer_type/:peer_id", e.listMessages)

		chatGroup.PUT("/:peer_type/:peer_id/:id", e.updateMessage)
		chatGroup.PUT("/:peer_type/:peer_id/read", e.readMessages)
		chatGroup.PUT("/:peer_type/:peer_id/typing", e.updateTyping)
	}
}

// PeerParams common struct for peer_type, peer_id
type PeerParams struct {
	PeerID   string `uri:"peer_id" binding:"uuid4"`
	PeerType string `uri:"peer_type" binding:"oneof=user group"`
}

type CreateMessageBody struct {
	Text string `json:"text" binding:"min=1,max=1000"`
}

type CreateMessageResponse struct {
	MessageID string `json:"message_id"`
}

func (e *HttpEndpoint) createMessage(ctx *gin.Context) {
	var (
		params PeerParams
		body   CreateMessageBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	message, err := e.service.CreateMessage(ctx, userID, params.PeerID, params.PeerType, body.Text)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, CreateMessageResponse{
		MessageID: message.ID,
	})
}

type UpdateTypingBody struct {
	Typing bool `json:"typing"`
}

func (e *HttpEndpoint) updateTyping(ctx *gin.Context) {
	var (
		params PeerParams
		body   UpdateTypingBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.UpdateTyping(ctx, userID, params.PeerID, params.PeerType, body.Typing); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

type ListMessagesParams struct {
	PeerID   string `uri:"peer_id" binding:"omitempty,uuid4"`
	PeerType string `uri:"peer_type" binding:"omitempty,oneof=user group"`
}

type ListMessagesQuery struct {
	Offset int64 `form:"offset" binding:"min=0"`
	Limit  int64 `form:"limit,default=100" binding:"min=1"`
}

type ListMessagesResponse struct {
	Items []chatdomain.MessageDTO `json:"items"`
}

func (e *HttpEndpoint) listMessages(ctx *gin.Context) {
	var (
		params ListMessagesParams
		query  ListMessagesQuery
	)

	if !transport.HttpBind(ctx, &params, nil, &query) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	var (
		messages []chatdomain.Message
		err      error
	)

	if params.PeerID != "" {
		messages, err = e.service.ListMessages(
			ctx,
			userID,
			params.PeerID,
			params.PeerType,
			query.Offset,
			query.Limit,
		)
		if err != nil {
			transport.HttpHandleError(ctx, err)

			return
		}
	} else {
		messages, err = e.service.ListRecentMessages(
			ctx,
			userID,
			params.PeerType,
			query.Offset,
			query.Limit,
		)
		if err != nil {
			transport.HttpHandleError(ctx, err)

			return
		}
	}

	items := lo.Map(messages, func(message chatdomain.Message, _ int) chatdomain.MessageDTO {
		return message.DTO()
	})

	ctx.JSON(http.StatusOK, ListMessagesResponse{
		Items: items,
	})
}

type UpdateMessageParams struct {
	PeerParams

	ID string `uri:"id" binding:"uuid4"`
}

type UpdateMessageBody struct {
	Text string `json:"text" binding:"min=1,max=1000"`
}

func (e *HttpEndpoint) updateMessage(ctx *gin.Context) {
	var (
		params UpdateMessageParams
		body   UpdateMessageBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	_, err := e.service.UpdateMessage(
		ctx,
		userID,
		params.PeerID,
		params.PeerType,
		params.ID,
		body.Text,
	)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

func (e *HttpEndpoint) readMessages(ctx *gin.Context) {
	var params PeerParams

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.ReadMessages(ctx, userID, params.PeerID, params.PeerType); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
