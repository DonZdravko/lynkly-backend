package common

import (
	"context"
	"net/http"
)

type ContextPair struct {
	Key   string
	Value interface{}
}

func ContextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func ContextSet(r *http.Request, contextValues ...ContextPair) *http.Request {
	if len(contextValues) == 0 {
		return r
	}

	ctx := r.Context()
	for _, value := range contextValues {
		ctx = context.WithValue(ctx, value.Key, value.Value)
	}

	return r.WithContext(ctx)
}

func SetPairsToContext(ctx context.Context, contextValues ...ContextPair) context.Context {
	for _, value := range contextValues {
		ctx = context.WithValue(ctx, value.Key, value.Value)
	}
	return ctx
}
