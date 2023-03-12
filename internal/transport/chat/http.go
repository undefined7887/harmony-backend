package chattransport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/undefined7887/harmony-backend/internal/domain"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	chatservice "github.com/undefined7887/harmony-backend/internal/service/chat"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/transport"
	authtransport "github.com/undefined7887/harmony-backend/internal/transport/auth"
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

func (e *HttpEndpoint) createMessage(ctx *gin.Context) {
	var (
		params chatdomain.PeerParams
		body   chatdomain.CreateMessageRequestBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	message, err := e.service.CreateMessage(
		ctx,
		userID,
		params.PeerID,
		params.PeerType,
		body.Text,
	)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, chatdomain.CreateMessageResponse{
		MessageID: message.ID,
	})
}

func (e *HttpEndpoint) getMessage(ctx *gin.Context) {
	var params domain.IdParam

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	message, err := e.service.GetMessage(ctx, userID, params.ID)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, chatdomain.GetMessageResponse{
		MessageDTO: message,
	})
}

func (e *HttpEndpoint) listMessages(ctx *gin.Context) {
	var (
		params chatdomain.PeerParams
		query  domain.PaginationQuery
	)

	if !transport.HttpBind(ctx, &params, nil, &query) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	messages, err := e.service.ListMessages(
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

	ctx.JSON(http.StatusOK, chatdomain.ListMessagesResponse{
		Items: messages,
	})
}

func (e *HttpEndpoint) updateMessage(ctx *gin.Context) {
	var (
		params domain.IdParam
		body   chatdomain.UpdateMessageRequestBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if _, err := e.service.UpdateMessage(
		ctx,
		userID,
		params.ID,
		body.Text,
	); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

func (e *HttpEndpoint) listChats(ctx *gin.Context) {
	var query chatdomain.ListChatsRequestQuery

	if !transport.HttpBindQuery(ctx, &query) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	chats, err := e.service.ListChats(
		ctx,
		userID,
		query.PeerType,
		query.Offset,
		query.Limit,
	)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, chatdomain.ListChatsResponse{
		Items: chats,
	})
}

func (e *HttpEndpoint) updateChatRead(ctx *gin.Context) {
	var params chatdomain.PeerParams

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

func (e *HttpEndpoint) updateChatTyping(ctx *gin.Context) {
	var (
		params chatdomain.PeerParams
		body   chatdomain.UpdateChatTypingRequestBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.UpdateChatTyping(
		ctx,
		userID,
		params.PeerID,
		params.PeerType,
		body.Typing,
	); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
