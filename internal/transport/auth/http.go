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

type SignUpRequest struct {
	Idtoken  string `json:"idtoken" binding:"jwt"`
	Nickname string `json:"nickname" binding:"nickname"`
}

func (e *HttpEndpoint) googleSignUp(c *gin.Context) {
	request := SignUpRequest{}

	if !transport.HttpBindJSON(c, &request) {
		return
	}

	token, err := e.service.GoogleSignUp(c, request.Idtoken, request.Nickname)
	if err != nil {
		transport.HttpHandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}

type SignInRequest struct {
	Idtoken string `json:"idtoken" binding:"jwt"`
}

func (e *HttpEndpoint) googleSignIn(c *gin.Context) {
	request := SignInRequest{}

	if err := c.BindJSON(&request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	token, err := e.service.GoogleSignIn(c, request.Idtoken)
	if err != nil {
		transport.HttpHandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}
