package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

const (
	keyParam   = "key"
	indexParam = "index"
	hkeyParam  = "hkey"
)

type ctxKey int

const (
	ctxKeyName ctxKey = iota
	ctxSetBody
	ctxHSetBody
	ctxRPushBody
	ctxIndex
	ctxHKeyName
)

// RequireKeyName middleware checks that 'key' parameter is set.
func RequireKeyName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, keyParam)
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "key is required"})

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

// RequireHKeyName middleware checks that 'hkey' parameter is set.
func RequireHKeyName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, hkeyParam)
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "hkey is required"})

			return
		}

		ctx := context.WithValue(r.Context(), ctxHKeyName, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetKeyName retrieves hash map key name value from context.
func GetHKeyName(ctx context.Context) string {
	v, ok := ctx.Value(ctxHKeyName).(string)
	if !ok {
		return ""
	}

	return v
}

// RequireIndex middleware checks that 'index' parameter is set.
func RequireIndex(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index := chi.URLParam(r, indexParam)
		if index == "" {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "index is required"})

			return
		}

		// Validate index
		v, err := strconv.Atoi(index)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "index is invalid"})

			return
		}

		// FIXME: current implementation of lindex does not allow negative indexes
		//       fix when available
		if v < 0 {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "index can't be negative"})

			return
		}

		ctx := context.WithValue(r.Context(), ctxIndex, v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetIndex retrieves index value from context.
func GetIndex(ctx context.Context) int {
	v, ok := ctx.Value(ctxIndex).(int)
	if !ok {
		return 0
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

// HSetRequestBody represents hset request body.
type HSetRequestBody struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
	TTL   int                    `json:"ttl"`
}

func (b *HSetRequestBody) IsValid() bool {
	return b.Key != "" && b.Value != nil
}

// RequireHSetParams validates request body for 'hset' operation.
func RequireHSetParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hsetBody := HSetRequestBody{}
		err := json.NewDecoder(r.Body).Decode(&hsetBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "hset body is invalid"})

			return
		}

		// Validate set body
		if !hsetBody.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": "hset body is invalid"})

			return
		}

		ctx = context.WithValue(ctx, ctxHSetBody, hsetBody)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetHSetBody retrieves set body from context.
func GetHSetBody(ctx context.Context) *HSetRequestBody {
	v, ok := ctx.Value(ctxHSetBody).(HSetRequestBody)
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
