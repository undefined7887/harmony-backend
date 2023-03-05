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
		// Message creation
		chatGroup.POST("/:peer_type/:peer_id", e.createMessage)
		chatGroup.POST("/:peer_type/:peer_id/typing", e.updateTyping)

		// Message listing
		chatGroup.GET("", e.listMessages)
		chatGroup.GET("/:peer_type", e.listMessages)
		chatGroup.GET("/:peer_type/:peer_id", e.listMessages)

		// Message updating
		chatGroup.PUT("/message/:id", e.updateMessage)

		// Message deleting
		chatGroup.DELETE("/message/:id", e.deleteMessage)
	}
}

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

type PeerParamsOmitempty struct {
	PeerID   string `uri:"peer_id" binding:"omitempty,uuid4"`
	PeerType string `uri:"peer_type" binding:"omitempty,oneof=user group"`
}

type PaginationQuery struct {
	Offset int64 `form:"offset" binding:"min=0"`
	Limit  int64 `form:"limit,default=100" binding:"min=1"`
}

type ListMessagesResponse struct {
	Items []chatdomain.MessageDTO `json:"items"`
}

func (e *HttpEndpoint) listMessages(ctx *gin.Context) {
	var (
		params PeerParamsOmitempty
		query  PaginationQuery
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

type MessageParamID struct {
	ID string `uri:"id" binding:"uuid4"`
}

type UpdateMessageBody struct {
	Text string `json:"text" binding:"min=1,max=1000"`
}

func (e *HttpEndpoint) updateMessage(ctx *gin.Context) {
	var (
		params MessageParamID
		body   UpdateMessageBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if _, err := e.service.UpdateMessage(ctx, userID, params.ID, body.Text); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

func (e *HttpEndpoint) deleteMessage(ctx *gin.Context) {
	params := MessageParamID{}

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	//userID := authtransport.GetClaims(ctx).Subject
	//
	//if _, err := e.service.UpdateMessage(ctx, userID, params.ID, request.Text); err != nil {
	//	transport.HttpHandleError(ctx, err)
	//
	//	return
	//}

	ctx.Status(http.StatusNotImplemented)
}
