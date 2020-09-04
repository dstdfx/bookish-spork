package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	keyParam = "key"
)

type ctxKey int

const (
	ctxKeyName ctxKey = iota
	ctxSetBody
	ctxRPushBody
)

// RequireKeyName middleware checks that 'key' parameter is set.
func RequireKeyName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, keyParam)
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, "key is required")

			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyName, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetKeyName retrieves key name value from context.
func GetKeyName(ctx context.Context) string {
	v, ok := ctx.Value(ctxKeyName).(string)
	if !ok {
		return ""
	}

	return v
}

// SetRequestBody represents set request body.
type SetRequestBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

func (b *SetRequestBody) IsValid() bool {
	return b.Key != "" && b.Value != nil
}

// RequireSetParams validates request body for 'set' operation.
func RequireSetParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		setBody := SetRequestBody{}
		err := json.NewDecoder(r.Body).Decode(&setBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "set body is invalid"})

			return
		}

		// Validate set body
		if !setBody.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "set body is invalid"})

			return
		}

		ctx = context.WithValue(ctx, ctxSetBody, setBody)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetSetBody retrieves set body from context.
func GetSetBody(ctx context.Context) *SetRequestBody {
	v, ok := ctx.Value(ctxSetBody).(SetRequestBody)
	if !ok {
		return nil
	}

	return &v
}

// RPushRequestBody represents rpush request body.
type RPushRequestBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

func (b *RPushRequestBody) IsValid() bool {
	return b.Key != "" && b.Value != nil
}

// RequireRPushParams validates request body for 'rpush' operation.
func RequireRPushParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rpush := RPushRequestBody{}
		err := json.NewDecoder(r.Body).Decode(&rpush)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "rpush body is invalid"})

			return
		}

		// Validate set body
		if !rpush.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "rpush body is invalid"})

			return
		}

		ctx = context.WithValue(ctx, ctxRPushBody, rpush)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRPushBody retrieves set body from context.
func GetRPushBody(ctx context.Context) *RPushRequestBody {
	v, ok := ctx.Value(ctxRPushBody).(RPushRequestBody)
	if !ok {
		return nil
	}

	return &v
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the Content-Type as application/json.
// It will call http.Error in case of failures.
func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
