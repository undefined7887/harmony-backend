package calltransport

import (
	"github.com/gin-gonic/gin"
	"github.com/undefined7887/harmony-backend/internal/domain"
	calldomain "github.com/undefined7887/harmony-backend/internal/domain/call"
	callservice "github.com/undefined7887/harmony-backend/internal/service/call"
	"github.com/undefined7887/harmony-backend/internal/transport"
	authtransport "github.com/undefined7887/harmony-backend/internal/transport/auth"
	"net/http"
)

type HttpEndpoint struct {
	service *callservice.Service
}

func NewHttpEndpoint(service *callservice.Service) transport.HttpEndpoint {
	return &HttpEndpoint{
		service: service,
	}
}

func (e *HttpEndpoint) Register(group *gin.RouterGroup) {
	callGroup := group.Group("/call")
	{
		callGroup.GET("", e.getCall)
		callGroup.PUT("/:id/status", e.updateCallStatus)
		callGroup.PUT("/:id/data", e.proxyCallData)

		callGroup.POST("/user/:peer_id", e.createCall)
	}
}

func (e *HttpEndpoint) createCall(ctx *gin.Context) {
	var (
		params calldomain.PeerParams
		body   calldomain.CreateCallRequestBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	callID, err := e.service.CreateCall(ctx, userID, params.PeerID, body.WebRtcOffer)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, calldomain.CreateCallResponse{
		CallID: callID,
	})
}

func (e *HttpEndpoint) getCall(ctx *gin.Context) {
	userID := authtransport.GetClaims(ctx).Subject

	call, err := e.service.GetCall(ctx, userID)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, call)
}

func (e *HttpEndpoint) updateCallStatus(ctx *gin.Context) {
	var (
		params domain.IdParam
		body   calldomain.UpdateCallRequestBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.UpdateCallStatus(ctx, userID, params.ID, body.Accept, body.WebRtcAnswer); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

func (e *HttpEndpoint) proxyCallData(ctx *gin.Context) {
	var (
		params domain.IdParam
		body   calldomain.ProxyCallDataRequestBody
	)

	if !transport.HttpBind(ctx, &params, &body, nil) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.ProxyCallData(ctx, userID, params.ID, body.WebRtcCandidate); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
