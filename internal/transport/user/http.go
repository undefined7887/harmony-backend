package usertransport

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	userservice "github.com/undefined7887/harmony-backend/internal/service/user"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"github.com/undefined7887/harmony-backend/internal/transport"
	authtransport "github.com/undefined7887/harmony-backend/internal/transport/auth"
)

const (
	SelfKeyword = "self"
)

type HttpEndpoint struct {
	service    *userservice.Service
	jwtService *jwtservice.Service
}

func NewHttpEndpoint(service *userservice.Service, jwtService *jwtservice.Service) transport.HttpEndpoint {
	return &HttpEndpoint{
		service:    service,
		jwtService: jwtService,
	}
}

func (e *HttpEndpoint) Register(group *gin.RouterGroup) {
	userGroup := group.
		Group("/user").
		Use(authtransport.NewHttpAuthMiddleware(e.jwtService))
	{
		userGroup.GET("/:id", e.getUser)
		userGroup.GET("/search", e.searchUser)
		userGroup.PUT("/status", e.updateUserStatus)
	}

	centrifugoGroup := group.
		Group("/user/centrifugo")
	{
		centrifugoGroup.POST("/connect", e.centrifugoConnect)
		centrifugoGroup.POST("/refresh", e.centrifugoRefresh)
	}
}

func (e *HttpEndpoint) getUser(ctx *gin.Context) {
	var params userdomain.GetUserRequestParams

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	if params.ID == SelfKeyword {
		params.ID = authtransport.GetClaims(ctx).Subject
	}

	user, err := e.service.GetUser(ctx, params.ID)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, userdomain.GetUserResponse{
		UserDTO: user,
	})
}

func (e *HttpEndpoint) searchUser(ctx *gin.Context) {
	var params userdomain.GetUserByNicknameRequestQuery

	if !transport.HttpBindQuery(ctx, &params) {
		return
	}

	user, err := e.service.SearchUser(ctx, params.Nickname)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, userdomain.GetUserResponse{
		UserDTO: user,
	})
}

func (e *HttpEndpoint) updateUserStatus(ctx *gin.Context) {
	var body userdomain.UpdateUserStatusRequestBody

	if !transport.HttpBindJSON(ctx, &body) {
		return
	}

	userID := authtransport.GetClaims(ctx).Subject

	if err := e.service.UpdateStatus(ctx, userID, body.Status, false); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

// Centrifugo events

func (e *HttpEndpoint) centrifugoConnect(ctx *gin.Context) {
	claims, err := authtransport.ExtractAndValidateToken(ctx, e.jwtService)
	if err != nil {
		ctx.JSON(http.StatusOK, userdomain.CentrifugoUnauthorizedResponse)

		return
	}

	if err := e.service.UpdateStatus(ctx, claims.Subject, userdomain.StatusOnline, true); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(
		http.StatusOK,
		centrifugo.NewResponse(userdomain.CentrifugoConnectResponse{
			User:     claims.Subject,
			ExpireAt: time.Now().Add(userdomain.UserPingInterval).Unix(),
		}),
	)
}

func (e *HttpEndpoint) centrifugoRefresh(ctx *gin.Context) {
	claims, err := authtransport.ExtractAndValidateToken(ctx, e.jwtService)
	if err != nil {
		ctx.JSON(http.StatusOK, userdomain.CentrifugoUnauthorizedResponse)

		return
	}

	if err := e.service.UpdateStatus(ctx, claims.Subject, userdomain.StatusOnline, true); err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(
		http.StatusOK,
		centrifugo.NewResponse(userdomain.CentrifugoRefreshResponse{
			ExpireAt: time.Now().Add(userdomain.UserPingInterval).Unix(),
		}),
	)
}
