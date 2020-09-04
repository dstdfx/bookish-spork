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
