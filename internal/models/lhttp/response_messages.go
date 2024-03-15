package lhttp

import (
	"errors"
	"net/http"
)

var (
	ErrNilError       = errors.New("error cannot be nil")
	ErrEmptyMessage   = errors.New("message cannot be empty")
	ErrNilJSONPayload = errors.New("json payload cannot be nil in favour of using NoContent instead")
)

type (
	PartialOK struct {
		PartialSuccess
	}

	PartialSuccess struct {
		response *HttpResponse
	}

	PartialFail struct {
		response *HttpResponse
	}

	PartialRedirect struct {
		response *HttpResponse
	}
)

// Info is used for informational responses where only headers are being passed back.
// Best example would be a HEAD request.
func (p *PartialOK) Info(contentType string, headers map[string]string) *HttpResponse {
	resp := p.response

	if contentType != "" {
		resp.contentType = contentType
	}

	resp.headers = headers
	return resp
}

// WithJSON adds a payload to the response and prepares it to be serialized as a JSON result.
// Panics if payload is nil.
func (p *PartialSuccess) WithJSON(payload interface{}) *HttpResponse {
	if payload == nil {
		panic(ErrNilJSONPayload)
	}

	p.response.payload = payload
	p.response.payloadType = jsonType
	p.response.contentType = ContentAppJSON

	return p.response
}

// WithText adds a text payload to the response.
func (p *PartialSuccess) WithText(text string) *HttpResponse {
	p.response.payload = text
	p.response.payloadType = textType
	p.response.contentType = ContentTextPlain

	return p.response
}

func (p *PartialSuccess) Etag() *PartialSuccess {
	p.response.etag = true

	return p
}

// FromTrustedError adds an error with trusted content to the response. This error will be displayed to the
// client, so it should NOT contain any internal information. Panics if error is nil.
func (p *PartialFail) FromTrustedError(err error) *HttpResponse {
	if err == nil {
		panic(ErrNilError)
	}

	if err.Error() == "" {
		panic(ErrEmptyMessage)
	}

	p.response.payload = err.Error()
	return p.response
}

// FromTrustedMessage adds an error message with trusted content to the response. This message will be displayed to the
// client, so it should NOT contain any internal information. Panics if message is empty.
func (p *PartialFail) FromTrustedMessage(publicMessage string) *HttpResponse {
	if publicMessage == "" {
		panic(ErrEmptyMessage)
	}

	p.response.payload = publicMessage
	return p.response
}

// WithJSON adds a payload to the response and prepares it to be serialized as a JSON result.
// Panics if payload is nil.
func (p *PartialFail) WithJSON(payload interface{}) *HttpResponse {
	if payload == nil {
		panic(ErrNilJSONPayload)
	}

	p.response.payload = payload
	p.response.payloadType = jsonType
	p.response.contentType = ContentAppJSON

	return p.response
}

func (p *PartialRedirect) Permanent(url string) *HttpResponse {
	p.response.statusCode = http.StatusPermanentRedirect
	p.response.payload = url
	p.response.contentType = ContentTextHTML

	return p.response
}

func (p *PartialRedirect) Temporary(url string) *HttpResponse {
	p.response.statusCode = http.StatusTemporaryRedirect
	p.response.payload = url
	p.response.contentType = ContentTextHTML

	return p.response
}

func (p *PartialRedirect) Found(url string) *HttpResponse {
	p.response.statusCode = http.StatusFound
	p.response.payload = url
	p.response.contentType = ContentTextHTML

	return p.response
}

//region 2xx

// OK returns a PartialSuccess. To use it as a response you should choose the payload next.
// In the cases of empty payloads use NoContent instead.
func OK() *PartialOK {
	return &PartialOK{PartialSuccess{newResponse(http.StatusOK)}}
}

// Created returns a PartialSuccess. To use it as a response you should choose the payload next.
// In the cases of empty payloads use NoContent instead.
func Created() *PartialSuccess {
	return &PartialSuccess{newResponse(http.StatusCreated)}
}

// Accepted returns a PartialSuccess. To use it as a response you should choose the payload next.
// In the cases of empty payloads use NoContent instead.
func Accepted() *PartialSuccess {
	return &PartialSuccess{newResponse(http.StatusAccepted)}
}

// NoContent returns a response with an empty payload.
func NoContent() *HttpResponse {
	return newResponse(http.StatusNoContent)
}

/* feel free to add 2xx responses above this line */

//endregion 2xx

//region 3xx

func NotModified() *HttpResponse {
	return newResponse(http.StatusNotModified)
}

func Redirect() *PartialRedirect {
	resp := &HttpResponse{
		isSuccessful: true,
		payloadType:  redirectType,
	}

	return &PartialRedirect{
		response: resp,
	}
}

//endregion 3xx

//region 4xx

// BadRequest returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func BadRequest() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusBadRequest,
	)}
}

// Unauthorized returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func Unauthorized() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusUnauthorized,
	)}
}

// PaymentRequired returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func PaymentRequired() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusPaymentRequired,
	)}
}

// Forbidden returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func Forbidden() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusForbidden,
	)}
}

// NotFound returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func NotFound() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusNotFound,
	)}
}

// MethodNotAllowed returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func MethodNotAllowed() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusMethodNotAllowed,
	)}
}

// Conflict returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func Conflict() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusConflict,
	)}
}

// NotAcceptable returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func NotAcceptable() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusNotAcceptable,
	)}
}

// Locked returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func Locked() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusLocked,
	)}
}

func RequestEntityTooLarge() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusRequestEntityTooLarge,
	)}
}

func StatusGone() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusGone,
	)}
}

func PreconditionFailed() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusPreconditionFailed,
	)}
}

/* feel free to add 4xx responses above this line */

//endregion 4xx

//region 5xx

// InternalServerError returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func InternalServerError() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusInternalServerError,
	)}
}

// NotImplemented returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func NotImplemented() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusNotImplemented,
	)}
}

// Unavailable returns a PartialFail. To use it as a response you need to select either FromTrustedError
// or FromTrustedMessage and provide a user-friendly info in both cases.
func Unavailable() *PartialFail {
	return &PartialFail{newErrorResponse(
		http.StatusServiceUnavailable,
	)}
}

/* feel free to add 5xx responses above this line */

//endregion 5xx

// CustomError is NOT intended to be used freely, but rather strictly and accompanied by
// a very good reason. It is unfortunately necessary for dealing with proxying of responses
// for external requests.
//
// Status code should be in the range [400-599]. It will otherwise default to 500.
func CustomError(status int) *PartialFail {
	if status < 400 || status > 599 {
		status = 500
	}

	return &PartialFail{newErrorResponse(status)}
}
