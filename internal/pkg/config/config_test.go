package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigInitFromStringValues(t *testing.T) {
	configString := `
log:
  file: "/var/log/test/test.log"
  use_stdout: true
  debug: true
public_api:
  server_address: localhost
  server_port: 63100
  read_timeout: 15
  write_timeout: 20
  idle_timeout: 30
service_api:
  server_address: localhost
  server_port: 63101
  read_timeout: 15
  write_timeout: 20
  idle_timeout: 30
cache:
  eviction_interval: 30
`

	expected := &AppConfig{
		Log: LogConfig{
			File:      "/var/log/test/test.log",
			UseStdout: true,
			Debug:     true,
		},
		PublicAPI: PublicAPIServerConfig{
			ServerAddress: "localhost",
			ServerPort:    63100,
			ReadTimeout:   15,
			WriteTimeout:  20,
			IdleTimeout:   30,
		},
		ServiceAPI: ServiceAPIServerConfig{
			ServerAddress: "localhost",
			ServerPort:    63101,
			ReadTimeout:   15,
			WriteTimeout:  20,
			IdleTimeout:   30,
		},
		Cache: CacheConfig{EvictionInterval: 30},
	}

	err := initFromString([]byte(configString))

	assert.Empty(t, err)
	assert.Equal(t, expected, Config)
}

func TestConfigInitFromStringDefaultValues(t *testing.T) {
	configString := ""

	expected := &AppConfig{
		PublicAPI: PublicAPIServerConfig{
			ServerAddress: "127.0.0.1",
			ServerPort:    63100,
			ReadTimeout:   60,
			WriteTimeout:  120,
			IdleTimeout:   240,
		},
		ServiceAPI: ServiceAPIServerConfig{
			ServerAddress: "127.0.0.1",
			ServerPort:    63101,
			ReadTimeout:   60,
			WriteTimeout:  120,
			IdleTimeout:   240,
		},
		Cache: CacheConfig{EvictionInterval: defaultEvictionInterval},
	}

	err := initFromString([]byte(configString))

	assert.Empty(t, err)
	assert.Equal(t, expected, Config)
}

func TestCheckConfigErr(t *testing.T) {
	Config = nil

	err := CheckConfig()

	assert.EqualError(t, err, errGlobalConfig)
}

func TestCheckConfig(t *testing.T) {
	Config = &AppConfig{Log: LogConfig{Debug: true}}

	err := CheckConfig()

	assert.NoError(t, err)
}
