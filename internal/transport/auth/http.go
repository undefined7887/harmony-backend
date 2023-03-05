package authtransport

import (
	"github.com/gin-gonic/gin"
	"github.com/undefined7887/harmony-backend/internal/service/auth"
	"github.com/undefined7887/harmony-backend/internal/transport"
	"net/http"
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

// AuthResponse is a base response for authentication endpoints
type AuthResponse struct {
	Token string `json:"token"`
}

type SignUpBody struct {
	Idtoken  string `json:"idtoken" binding:"jwt"`
	Nickname string `json:"nickname" binding:"nickname"`
}

func (e *HttpEndpoint) googleSignUp(ctx *gin.Context) {
	var body SignUpBody

	if !transport.HttpBindJSON(ctx, &body) {
		return
	}

	token, err := e.service.GoogleSignUp(ctx, body.Idtoken, body.Nickname)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}

type SignInBody struct {
	Idtoken string `json:"idtoken" binding:"jwt"`
}

func (e *HttpEndpoint) googleSignIn(ctx *gin.Context) {
	var body SignInBody

	if !transport.HttpBindJSON(ctx, &body) {
		return
	}

	token, err := e.service.GoogleSignIn(ctx, body.Idtoken)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}
