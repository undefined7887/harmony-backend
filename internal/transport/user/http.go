package usertransport

import (
	"github.com/gin-gonic/gin"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	userservice "github.com/undefined7887/harmony-backend/internal/service/user"
	"github.com/undefined7887/harmony-backend/internal/transport"
	authtransport "github.com/undefined7887/harmony-backend/internal/transport/auth"
	"net/http"
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
		userGroup.GET("/:id", e.read)
		userGroup.GET("/nickname", e.readByNickname)
	}
}

type ReadParams struct {
	ID string `uri:"id" binding:"uuid4|eq=self"`
}

func (e *HttpEndpoint) read(ctx *gin.Context) {
	var params ReadParams

	if !transport.HttpBindURI(ctx, &params) {
		return
	}

	if params.ID == SelfKeyword {
		params.ID = authtransport.GetClaims(ctx).Subject
	}

	user, err := e.service.Read(ctx, params.ID)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, user.DTO())
}

type ReadByNicknameBody struct {
	Nickname string `json:"nickname" binding:"nickname-extended"`
}

func (e *HttpEndpoint) readByNickname(ctx *gin.Context) {
	var params ReadByNicknameBody

	if !transport.HttpBindJSON(ctx, &params) {
		return
	}

	user, err := e.service.ReadByNickname(ctx, params.Nickname)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, user.DTO())
}
