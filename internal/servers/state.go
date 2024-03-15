package servers

import (
	"lynkly-backend/internal/logging"
	"lynkly-backend/internal/routers"
)

type State struct {
	Routers routers.RouteVersions
	// TODO: [Zdravko Donev] Add JWT middleware
	//Middleware *jwtmiddleware.JWTMiddleware
}

type ServerParams struct {
	Logger     logging.Logger
	ServiceUrl string
}
