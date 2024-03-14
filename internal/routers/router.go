package routers

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"lynkly-backend/internal/logging"
	"lynkly-backend/internal/models/lhttp"
	"net/http"
	"runtime/debug"
)

type Router struct {
	router *mux.Router
	logger logging.Logger
	// TODO: [Zdravko Donev] Add jwt middleware
	//jwt     *jwtmiddleware.JWTMiddleware
	options []MiddleWareOptions
}

func NewRouter(router *mux.Router, routerParams *RouterParams, options ...MiddleWareOptions) *Router {
	return &Router{
		router:  router,
		logger:  routerParams.Logger,
		options: options,
	}
}

func (tr *Router) HandleFunc(method, url string, handler RouteHandlerFunc) {
	tr.logger.Debug("Router.HandleFunc - Registering handler: ", method, " ", url)
	tr.handle(method, url, negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			tr.logger.Debug("Router.HandleFunc: ", r.URL.Path)
			next(w, r)
		}),
		tr.Wrap(handler)),
	)
}

func (tr *Router) Wrap(routeHandler RouteHandlerFunc) negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
		newRequest := r.WithContext(r.Context())

		resp := routeHandler(newRequest)

		if err := lhttp.Write(w, r, resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			tr.logger.WithRequest(newRequest).
				Error(err.Error())
		}

		if resp.IsSuccessful() == false {
			tr.logResponse(newRequest, resp)
		}
	})
}

func (tr *Router) logResponse(r *http.Request, resp *lhttp.HttpResponse) {
	if resp.Payload() == nil {
		return
	}

	errMsg, ok := resp.Payload().(string)
	if ok == false {
		errMsg = "error response did not hold the expected payload"
		tr.logger.WithRequest(r).Error(errMsg)
		return
	}

	if resp.StatusCode() >= 500 {
		tr.logger.WithRequest(r).Error(errMsg)
		return
	}

	// 4xx responses
	tr.logger.WithRequest(r).Warn(errMsg)
}

func (tr *Router) handle(method, url string, handler http.Handler) {
	tr.router.Handle(url, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr.logger.Debug("brej", r.URL.Path)
		defer func() {
			if rec := recover(); rec != nil {
				tr.logger.Panic("panic in router handler", string(debug.Stack()))

				resp := lhttp.InternalServerError().FromTrustedMessage(http.StatusText(http.StatusInternalServerError))
				_ = lhttp.Write(w, r, resp)
			}

			handler.ServeHTTP(w, r)
		}()
	})).Methods(method)
}
