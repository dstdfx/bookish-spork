package testutils

import (
	"os"
	"testing"

	"github.com/dstdfx/bookish-spork/internal/pkg/config"
)

const (
	accTestsEnabled    = "ACC_TESTS"
	accTestEnabledMsg  = "Acceptance tests suite is enabled"
	accTestDisabledMsg = `Acceptance tests suite is disabled, you can enable it with ACC_TESTS=1`
)

// InitTestConfig initializes application configuration for tests.
func InitTestConfig() {
	config.Config = &config.AppConfig{
		Log: config.LogConfig{
			UseStdout: true,
			Debug:     true,
		},
	}
}

// IsAccTestEnabled checks if aceptance tests are enabled.
func IsAccTestEnabled(t *testing.T) bool {
	if os.Getenv(accTestsEnabled) == "1" {
		t.Log(accTestEnabledMsg)

		return true
	}
	t.Log(accTestDisabledMsg)

	return false
}
