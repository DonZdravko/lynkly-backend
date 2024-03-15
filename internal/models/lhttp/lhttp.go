package lhttp

import (
	"encoding/json"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strings"
)

var ErrNotSupportedPayloadType = errors.New("unsupported payload type")

// custom jsoniter because hashing requires consistent content
// which is not possible when keys are not sorted
var etagJsoniter = jsoniter.Config{
	EscapeHTML:  true,
	SortMapKeys: true,
}.Froze()

type payloadType int

const (
	emptyType payloadType = iota
	jsonType
	textType
	redirectType
)

// HttpResponse holds the data for a http response. It is meant to be used by Write
// to populate its content to a http.ResponseWriter
type HttpResponse struct {
	statusCode   int
	isSuccessful bool
	contentType  string
	payload      interface{}
	payloadType  payloadType
	headers      map[string]string
	etag         bool
}

func (r *HttpResponse) StatusCode() int {
	return r.statusCode
}

func (r *HttpResponse) IsSuccessful() bool {
	return r.isSuccessful
}

func (r *HttpResponse) Payload() interface{} {
	return r.payload
}

// WithHeaders adds additional headers to the response. However, Content-Type and other headers already
// used by the responses will not be overridden.
//
// Calling this function again would override the previous results.
func (r *HttpResponse) WithHeaders(headers map[string]string) *HttpResponse {
	r.headers = headers
	return r
}

// WithHeader adds additional header to the response. However, Content-Type and other headers already
// used by the responses will not be overridden.
//
// Calling this function again would override the previous results.
func (r *HttpResponse) WithHeader(header http.Header) *HttpResponse {
	headers := make(map[string]string)
	for k, v := range header {
		headers[k] = strings.Join(v, " ")
	}

	r.headers = headers
	return r
}

func newResponse(statusCode int) *HttpResponse {
	return &HttpResponse{
		statusCode:   statusCode,
		isSuccessful: true,
		contentType:  ContentTextPlain,
		payload:      nil,
		payloadType:  emptyType,
		headers:      nil,
	}
}

func newErrorResponse(statusCode int) *HttpResponse {
	return &HttpResponse{
		statusCode:   statusCode,
		isSuccessful: false,
		contentType:  ContentTextPlain,
		payload:      nil,
		payloadType:  emptyType,
		headers:      nil,
	}
}

// headers
const (
	AuthorizationHeader       = "Authorization"
	AcceptHeader              = "Accept"
	ContentTypeHeader         = "Content-Type"
	ContentTypeOptions        = "X-Content-Type-Options"
	ContentTypeOptionsNoSniff = "nosniff"
	ContentAppJSON            = "application/json;charset=utf-8"
	ContentAppJSONPatch       = "application/json-patch+json"
	ContentFormURL            = "application/x-www-form-urlencoded"
	ContentTextPlain          = "text/plain;charset=utf-8"
	ContentAppPdf             = "application/pdf"
	ContentAppOctetStream     = "application/octet-stream"
	ContentTextYAML           = "text/yaml"
	ContentTextHTML           = "text/html;charset=utf-8"
)

// Write is writing the data of the response to the provided http.ResponseWriter.
func Write(w http.ResponseWriter, r *http.Request, response *HttpResponse) error {
	applyCustomHeaders(w, response.headers)
	w.Header().Set(ContentTypeHeader, response.contentType)

	if response.isSuccessful == false {
		return writeError(w, response)
	}

	return writePayload(w, r, response)
}

func redirect(response *HttpResponse, w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, response.payload.(string), response.statusCode)
	return nil
}

func writePayload(w http.ResponseWriter, r *http.Request, response *HttpResponse) error {
	if response.payloadType == redirectType {
		return redirect(response, w, r)
	}

	if response.etag {
		return writeEtagPayload(w, r, response)
	}

	//writing the header must be done after setting all headers and right before writing the payload
	w.WriteHeader(response.statusCode)

	switch response.payloadType {
	case textType:
		_, err := io.WriteString(w, response.payload.(string))
		return err
	case jsonType:
		return jsoniter.NewEncoder(w).Encode(&response.payload)
	case emptyType:
		return nil
	}

	return ErrNotSupportedPayloadType
}

func writeEtagPayload(w http.ResponseWriter, r *http.Request, response *HttpResponse) error {
	switch response.payloadType {
	case textType:
		res := []byte(response.payload.(string))
		if ok, err := handleEtag(w, r, res); ok {
			return err
		}

		w.WriteHeader(response.statusCode)
		_, err := w.Write(res)
		return err
	case jsonType:
		res, err := etagJsoniter.Marshal(&response.payload)
		if err != nil {
			return err
		}

		if ok, err := handleEtag(w, r, res); ok {
			return err
		}

		w.WriteHeader(response.statusCode)
		_, err = w.Write(res)
		return err
	case emptyType:
		return nil
	}

	return ErrNotSupportedPayloadType
}

func writeError(w http.ResponseWriter, response *HttpResponse) error {
	w.Header().Set(ContentTypeOptions, ContentTypeOptionsNoSniff)
	w.WriteHeader(response.statusCode)

	if response.payloadType == jsonType {
		return json.NewEncoder(w).Encode(response.payload)
	}

	_, err := io.WriteString(w, response.payload.(string))
	return err
}

func applyCustomHeaders(w http.ResponseWriter, headers map[string]string) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}
