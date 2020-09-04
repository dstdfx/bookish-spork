package http

import (
	"github.com/dstdfx/bookish-spork/internal/pkg/backend"
	v1 "github.com/dstdfx/bookish-spork/internal/pkg/http/v1"
	"github.com/go-chi/chi"
)

const (
	groupV1 = "/v1"
)

// InitAPIRouter configures HTTP router.
func InitAPIRouter(b *backend.Backend) chi.Router {
	r := chi.NewRouter()
	r.Route(groupV1, func(r chi.Router) {
		r.Mount("/", v1.Routes(b))
	})

	return r
}
