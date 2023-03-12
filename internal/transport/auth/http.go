package authtransport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authdomain "github.com/undefined7887/harmony-backend/internal/domain/auth"
	"github.com/undefined7887/harmony-backend/internal/service/auth"
	"github.com/undefined7887/harmony-backend/internal/transport"
)

type HttpEndpoint struct {
	service *authservice.Service
}

func NewHttpEndpoint(service *authservice.Service) transport.HttpEndpoint {
	return &HttpEndpoint{
		service: service,
	}
}

func (e *HttpEndpoint) Register(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/google/sign_up", e.googleSignUp)
		authGroup.POST("/google/sign_in", e.googleSignIn)
	}
}

func (e *HttpEndpoint) googleSignUp(ctx *gin.Context) {
	var body authdomain.SignUpRequestBody

	if !transport.HttpBindJSON(ctx, &body) {
		return
	}

	auth, err := e.service.GoogleSignUp(ctx, body.Idtoken, body.Nickname)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, authdomain.SignUpResponse{
		AuthDTO: auth,
	})
}

func (e *HttpEndpoint) googleSignIn(ctx *gin.Context) {
	var body authdomain.SignInRequestBody

	if !transport.HttpBindJSON(ctx, &body) {
		return
	}

	auth, err := e.service.GoogleSignIn(ctx, body.Idtoken)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, authdomain.SignInResponse{
		AuthDTO: auth,
	})
}
