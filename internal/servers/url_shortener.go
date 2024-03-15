package servers

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"lynkly-backend/internal/logging"
	"lynkly-backend/internal/models/lhttp"
	"lynkly-backend/internal/routers"
	"math/rand"
	"net/http"
	"net/url"
)

// TODO: [Zdravko Donev] Remove this after we get the DB
var shortURLs map[string]string

type UrlShortenerServer struct {
	hostPort   string
	handler    http.Handler
	logger     logging.Logger
	serviceUrl string
}

func NewUrlShortenerServer(port string, serverParams ServerParams) *UrlShortenerServer {
	shortURLs = make(map[string]string)
	muxRouter := mux.NewRouter().StrictSlash(false)
	state := &State{
		Routers: routers.RouteVersions{
			V1: routers.NewRouter(muxRouter.PathPrefix(routers.PathAPIV1).Subrouter(), &routers.RouterParams{
				Logger: serverParams.Logger,
			}),
		},
	}

	corsHandler := AddCorsOptions(muxRouter, cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowedHeaders: []string{"*"},
	})

	urlShortenerServer := &UrlShortenerServer{
		hostPort:   port,
		handler:    corsHandler,
		logger:     serverParams.Logger,
		serviceUrl: serverParams.ServiceUrl,
	}

	urlShortenerServer.registerApiHandlers(state)

	return urlShortenerServer
}

func (s *UrlShortenerServer) Run() error {
	s.logger.Info("Starting url shortener API", "port: "+s.hostPort)
	return http.ListenAndServe(s.hostPort, s.handler)
}

func (s *UrlShortenerServer) registerApiHandlers(state *State) {
	state.Routers.V1.HandleFunc(http.MethodGet, "/{shortURL}", s.RedirectHandler)
	state.Routers.V1.HandleFunc(http.MethodPost, "/shorten", s.ShortenHandler)
}

func (s *UrlShortenerServer) RedirectHandler(r *http.Request) *lhttp.HttpResponse {
	s.logger.Debug("UrlShortenerServer.RedirectHandler")
	shortURL, found := mux.Vars(r)["shortURL"]
	if found == false {
		s.logger.WithRequest(r).Error("Short URL not found in request")
		return lhttp.BadRequest().FromTrustedMessage("Short URL not found in request")
	}
	s.logger.Debug("Short URL found in request: ", shortURL)

	for key, value := range shortURLs {
		s.logger.Debug("shortURLs: ", key, " - ", value)
	}
	longURL, ok := shortURLs[shortURL]
	if !ok {
		return lhttp.NotFound().FromTrustedMessage(fmt.Sprintf("Short URL not found - %s", shortURL))
	}
	s.logger.Debug("Long URL found for short URL: ", longURL)

	return lhttp.Redirect().Temporary(longURL)
}

func (s *UrlShortenerServer) ShortenHandler(r *http.Request) *lhttp.HttpResponse {
	s.logger.Debug("UrlShortenerServer.ShortenHandler")
	s.logger.Info("Shortening URL")
	longURL := r.FormValue("url")
	if longURL == "" {
		return lhttp.BadRequest().FromTrustedMessage("Missing URL parameter")
	}
	// Check if longURL is a valid URL
	if parsedUrl, err := url.Parse(longURL); err != nil {
		return lhttp.BadRequest().FromTrustedMessage("Invalid URL - " + longURL + " - " + err.Error())
	} else if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return lhttp.BadRequest().FromTrustedMessage("Invalid URL - " + longURL + " - " + "URL is missing scheme or host")
	}

	shortCode := s.GenerateShortURL()
	// TODO: Store the shortURL and longURL in a database
	shortURLs[shortCode] = longURL
	shortURL := s.serviceUrl + routers.PathAPIV1 + "/" + shortCode

	s.logger.Info("Shortened URL: " + shortURL)
	return lhttp.OK().WithJSON(map[string]string{
		"shortUrl": shortURL,
	})
}

func (s *UrlShortenerServer) GenerateShortURL() string {
	// Generate a short URL logic here (you can use hashing algorithms or any other method)
	// For simplicity, we'll just take the first 6 characters of the URL
	shortCode := base64.StdEncoding.EncodeToString(Int63ToByteArray(rand.Uint64()))
	s.logger.Debug("Generated short code: ", shortCode)

	return shortCode
	//if len(longURL) <= 6 {
	//	return longURL
	//}
	//return longURL[:6]
}

func Int63ToByteArray(value uint64) []byte {
	var buffer [8]byte                           // Allocate a fixed-size byte array
	binary.BigEndian.PutUint64(buffer[:], value) // Use BigEndian or LittleEndian depending on your needs
	return buffer[:]
}
