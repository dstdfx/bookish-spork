package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/dstdfx/bookish-spork/internal/pkg/backend"
	"github.com/dstdfx/bookish-spork/internal/pkg/qqcache"
	"github.com/go-chi/chi"
)

// Routes initializes v1 handler.
func Routes(b *backend.Backend) http.Handler {
	r := chi.NewRouter()

	// GET /v1/get/<key>
	r.
		With(RequireKeyName).
		Get("/get/{key}", getHandler(b))

	// POST /v1/set
	r.
		With(RequireSetParams).
		Post("/set", setHandler(b))

	// GET /v1/keys
	r.Get("/keys", keysHandler(b))

	// DELETE /v1/remove/<key>
	r.
		With(RequireKeyName).
		Delete("/remove/{key}", removeHandler(b))

	// POST /v1/rpush
	r.
		With(RequireRPushParams).
		Post("/rpush", rpushHandler(b))

	// GET /v1/lindex/<key>/<index>
	r.
		With(RequireKeyName).
		With(RequireIndex).
		Get("/lindex/{key}/{index}", lindexHandler(b))

	// POST /v1/hset
	r.
		With(RequireHSetParams).
		Post("/hset", hsetHandler(b))

	// GET /v1/hget/<key>/<hkey>
	r.
		With(RequireKeyName).
		With(RequireHKeyName).
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
		// Get rpush body from router's context
		body := GetRPushBody(req.Context())

		err := b.Cache.RPush(body.Key, body.Value, time.Duration(body.TTL)*time.Second)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": err.Error()})

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func lindexHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get key from router's context
		key := GetKeyName(req.Context())
		index := GetIndex(req.Context())

		v, err := b.Cache.LIndex(key, index)
		if err != nil {
			if errors.Is(err, qqcache.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			}
			if errors.Is(err, qqcache.ErrWrongTypeIndex) {
				w.WriteHeader(http.StatusBadRequest)
				JSON(w, map[string]string{"error": err.Error()})
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		JSON(w, map[string]interface{}{"value": v})
	}
}

func hsetHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get hset body from router's context
		body := GetHSetBody(req.Context())

		err := b.Cache.HSet(body.Key, body.Value, time.Duration(body.TTL)*time.Second)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": err.Error()})

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func hgetHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get key from router's context
		key := GetKeyName(req.Context())
		hkey := GetHKeyName(req.Context())

		v, err := b.Cache.HGet(key, hkey)
		if err != nil {
			if errors.Is(err, qqcache.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			}
			if errors.Is(err, qqcache.ErrWrongTypeHGet) {
				w.WriteHeader(http.StatusBadRequest)
				JSON(w, map[string]string{"error": err.Error()})
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		JSON(w, map[string]interface{}{"value": v})
	}
}
