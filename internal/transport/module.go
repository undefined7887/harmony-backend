package transport

import (
	"go.uber.org/fx"

	_ "github.com/undefined7887/harmony-backend/internal/validation" // enabling validation
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(
		NewHttpServer,
		fx.ParamTags("", "", `group:"http_endpoints"`),
	)),
	fx.Invoke(NewHttpServerRunner),
)
