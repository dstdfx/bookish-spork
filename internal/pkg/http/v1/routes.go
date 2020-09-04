package v1

import (
	"net/http"
	"time"

	"github.com/dstdfx/bookish-spork/internal/pkg/backend"
	"github.com/go-chi/chi"
)

// Routes initializes v1 handler.
func Routes(b *backend.Backend) http.Handler {
	r := chi.NewRouter()

	// GET /v1/get/<key>
	r.With(RequireKeyName).
		Get("/get/{key}", getHandler(b))

	// POST /v1/set
	r.With(RequireSetParams).
		Post("/set", setHandler(b))

	// GET /v1/keys
	r.Get("/keys", keysHandler(b))

	// DELETE /v1/remove/<key>
	r.With(RequireKeyName).
		Delete("/remove/{key}", removeHandler(b))

	// POST /v1/rpush
	r.Post("/rpush", rpushHandler(b))

	// GET /v1/lindex/<key>/<index>
	r.With(RequireKeyName).
		Get("/lindex/{key}/{index}", lindexHandler(b))

	// POST /v1/hset
	r.Post("/hset", hsetHandler(b))

	// GET /v1/hget/<key>/<hkey>
	r.With(RequireKeyName).
		Get("/hget/{key}/{hkey}", hgetHandler(b))

	return r
}

func getHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get key from router's context
		key := GetKeyName(req.Context())

		// Get value from cache
		k, ok := b.Cache.Get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.WriteHeader(http.StatusOK)
		JSON(w, map[string]interface{}{"value": k})
	}
}

func setHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get set body from router's context
		body := GetSetBody(req.Context())

		// Set new entity
		b.Cache.Set(body.Key, body.Value, time.Duration(body.TTL)*time.Second)
		w.WriteHeader(http.StatusOK)
	}
}

func keysHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		JSON(w, map[string]interface{}{"keys": b.Cache.Keys()})
	}
}

func removeHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get key from router's context
		key := GetKeyName(req.Context())

		// Remove key from the cache
		b.Cache.Remove(key)
		w.WriteHeader(http.StatusNoContent)
	}
}

func rpushHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func lindexHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func hsetHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func hgetHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
