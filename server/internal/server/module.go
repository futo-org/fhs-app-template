package server

import (
	"go.uber.org/fx"
)

var Module = fx.Module("server",
	fx.Provide(NewServer),
	fx.Invoke(RegisterRoutes),
	fx.Invoke(RegisterWebUI),
	fx.Invoke(RunServer),
)
