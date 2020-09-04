package bookishspork

import (
	"os"
	"sync"
	"syscall"
	"testing"

	"github.com/dstdfx/bookish-spork/internal/pkg/config"
	"github.com/dstdfx/bookish-spork/internal/pkg/log"
	"github.com/dstdfx/bookish-spork/internal/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStartService(t *testing.T) {
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

	var wg sync.WaitGroup
	wg.Add(1)
	interrupt := make(chan os.Signal, 1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		assert.NoError(t, StartService(logger, StartOpts{Interrupt: interrupt}))
	}(&wg)

	// Send interrupt.
	interrupt <- syscall.SIGINT

	wg.Wait()
}
