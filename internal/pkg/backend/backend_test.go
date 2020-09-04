package backend

import (
	"testing"

	"github.com/dstdfx/bookish-spork/internal/pkg/config"
	"github.com/dstdfx/bookish-spork/internal/pkg/log"
	"github.com/dstdfx/bookish-spork/internal/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBackend(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	b := New(logger)
	defer b.Shutdown()
	assert.NotNil(t, b)
}
