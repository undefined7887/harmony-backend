package usertransport

import (
	"github.com/gin-gonic/gin"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	userservice "github.com/undefined7887/harmony-backend/internal/service/user"
	"github.com/undefined7887/harmony-backend/internal/transport"
	authtransport "github.com/undefined7887/harmony-backend/internal/transport/auth"
	httputil "github.com/undefined7887/harmony-backend/internal/util/http"
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
		userGroup.GET("/:id", e.get)
	}
}

type GetRequest struct {
	ID string `uri:"id" binding:"uuid|eq=self"`
}

func (r *GetRequest) Process(ctx *gin.Context) {
	if r.ID == SelfKeyword {
		r.ID = httputil.GetClaims(ctx).Subject
	}
}

type GetResponse struct {
	User *userdomain.UserDTO `json:"user"`
}

func (e *HttpEndpoint) get(ctx *gin.Context) {
	request := GetRequest{}

	if !transport.HttpBindURI(ctx, &request) {
		return
	}

	// Processing all possible values of id
	request.Process(ctx)

	user, err := e.service.Get(ctx, request.ID)
	if err != nil {
		transport.HttpHandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, GetResponse{
		User: user.DTO(),
	})
}
