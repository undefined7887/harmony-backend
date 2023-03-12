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
		chatGroup.GET("", e.listChats)

		chatGroup.POST("/:peer_type/:peer_id", e.createMessage)
		chatGroup.GET("/:peer_type/:peer_id", e.listMessages)
		chatGroup.PUT("/:peer_type/:peer_id/read", e.updateChatRead)
		chatGroup.PUT("/:peer_type/:peer_id/typing", e.updateChatTyping)

		chatGroup.GET("/message/:id", e.getMessage)
		chatGroup.PUT("/message/:id", e.updateMessage)
	}
}

// Common DTOs

type PeerParams struct {
	PeerID   string `uri:"peer_id" binding:"id"`
	PeerType string `uri:"peer_type" binding:"oneof=user group"`
}

type MessageIdParam struct {
	ID string `uri:"id" binding:"id"`
}

type PaginationQuery struct {
	Offset int64 `form:"offset" binding:"min=0"`
	Limit  int64 `form:"limit,default=100" binding:"min=1"`
}

// Handlers

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

func (e *HttpEndpoint) getMessage(ctx *gin.Context) {
	var params MessageIdParam

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	message, err := e.service.GetMessage(ctx, userID, params.ID)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, message.DTO())
}

type ListMessagesResponse struct {
	Items []chatdomain.MessageDTO `json:"items"`
}

func (e *HttpEndpoint) listMessages(ctx *gin.Context) {
	var (
		params PeerParams
		query  PaginationQuery
	)

	if !transport.HttpBind(ctx, &params, nil, &query) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	messages, err := e.service.ListMessages(ctx, userID, params.PeerID, params.PeerType, query.Offset, query.Limit)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, ListMessagesResponse{
		Items: lo.Map(messages, func(message chatdomain.Message, _ int) chatdomain.MessageDTO {
			return message.DTO()
		}),
	})
}

type UpdateMessageBody struct {
	Text string `json:"text" binding:"min=1,max=1000"`
}

func (e *HttpEndpoint) updateMessage(ctx *gin.Context) {
	var (
		params MessageIdParam
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

type ListChatsQuery struct {
	PaginationQuery

	PeerType string `form:"peer_type" binding:"omitempty,oneof=user group"`
}

type ListChatsResponse struct {
	Items []chatdomain.ChatDTO `json:"items"`
}

func (e *HttpEndpoint) listChats(ctx *gin.Context) {
	var query ListChatsQuery

	if !transport.HttpBindQuery(ctx, &query) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	chats, err := e.service.ListChats(ctx, userID, query.PeerType, query.Offset, query.Limit)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, ListChatsResponse{
		Items: lo.Map(chats, func(item chatdomain.Chat, _ int) chatdomain.ChatDTO {
			return item.DTO()
		}),
	})
}

func (e *HttpEndpoint) updateChatRead(ctx *gin.Context) {
	var params PeerParams

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.UpdateChatRead(ctx, userID, params.PeerID, params.PeerType); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

type UpdateTypingBody struct {
	Typing bool `json:"typing"`
}

func (e *HttpEndpoint) updateChatTyping(ctx *gin.Context) {
	var (
		params PeerParams
		body   UpdateTypingBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.UpdateChatTyping(ctx, userID, params.PeerID, params.PeerType, body.Typing); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
