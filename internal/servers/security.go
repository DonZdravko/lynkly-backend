package servers

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func AddCorsOptions(router *mux.Router, corsOptions cors.Options) http.Handler {
	corsOptions.MaxAge = 604800 // Tell the client to cache the CORS headers for 1 week

	return cors.New(corsOptions).Handler(router)
}
