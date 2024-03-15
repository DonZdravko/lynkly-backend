package routers

import (
	"lynkly-backend/internal/logging"
	"lynkly-backend/internal/models/lhttp"
	"net/http"
)

type RouterParams struct {
	Logger logging.Logger
}

const (
	PathAPIV1 = "/api/v1"
	PathAPIV2 = "/api/v2"
	PathAPIV3 = "/api/v3"
)

type RouteVersions struct {
	V1 *Router
}

type OptionType int
type MiddleWareOptions struct {
	Key   OptionType
	Value interface{}
}

type RouteHandlerFunc func(r *http.Request) *lhttp.HttpResponse
