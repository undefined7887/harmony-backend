package transport

import (
	_ "github.com/undefined7887/harmony-backend/internal/validation" // enabling validation
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(
		NewHttpServer,
		fx.ParamTags("", "", `group:"http_endpoints"`),
	)),
	fx.Invoke(NewHttpServerRunner),
)
