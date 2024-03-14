package lhttp

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/go-http-utils/fresh"
	"github.com/go-http-utils/headers"
	"net/http"
)

func handleEtag(w http.ResponseWriter, r *http.Request, res []byte) (bool, error) {
	if r == nil {
		return false, nil
	}

	resHeader := w.Header()

	etag, err := generateEtag(res)
	if err != nil {
		return false, err
	}

	resHeader.Set(headers.ETag, *etag)
	// max-age = 2 weeks
	resHeader.Set(headers.CacheControl, "max-age=1209600, private, no-cache")

	if fresh.IsFresh(r.Header, resHeader) {
		w.WriteHeader(http.StatusNotModified)
		_, err = w.Write(nil)
		return true, err
	}
	return false, nil
}

func generateEtag(res []byte) (*string, error) {
	hasher := sha1.New()
	if _, err := hasher.Write(res); err != nil {
		return nil, err
	}

	etag := fmt.Sprintf("W/%q", hex.EncodeToString(hasher.Sum(nil)))
	return &etag, nil
}
