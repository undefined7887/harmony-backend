package usertransport

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotated{
		Group:  "http_endpoints",
		Target: NewHttpEndpoint,
	},
)
