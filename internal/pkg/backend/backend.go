package backend

import (
	"time"

	"github.com/dstdfx/bookish-spork/internal/pkg/config"
	"github.com/dstdfx/bookish-spork/internal/pkg/qqcache"
	"go.uber.org/zap"
)

// Backend contains common application dependencies.
type Backend struct {
	Log   *zap.Logger
	Cache *qqcache.Cache
}

// New init new Backend instance.
func New(log *zap.Logger) *Backend {
	opts := qqcache.Opts{
		EvictionInterval: time.Duration(config.Config.Cache.EvictionInterval) * time.Second,
	}

	return &Backend{
		Log:   log,
		Cache: qqcache.New(opts),
	}
}

// Shutdown method closes all backend connections.
func (b *Backend) Shutdown() {
	b.Log.Debug("backend shutdown")
	b.Cache.Shutdown()
}
