package config

import (
	"errors"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

const (
	// errGlobalConfig represents error message in case of empty global config.
	errGlobalConfig = "global configuration is not initialized"

	defaultPublicAPIAddress = "127.0.0.1"
	defaultPublicAPIPort    = 63100

	defaultServiceAPIAddress = "127.0.0.1"
	defaultServiceAPIPort    = 63101

	defaultHTTPReadTimeout  = 60
	defaultHTTPWriteTimeout = 120
	defaultHTTPIdleTimeout  = 240
)

// Config is a global container for all configuration options.
var Config *AppConfig

// AppConfig contains all application parameters.
type AppConfig struct {
	Log        LogConfig              `yaml:"log"`
	PublicAPI  PublicAPIServerConfig  `yaml:"public_api"`
	ServiceAPI ServiceAPIServerConfig `yaml:"service_api"`
}

// LogConfig contains logger configuration.
type LogConfig struct {
	File      string `yaml:"file"`
	UseStdout bool   `yaml:"use_stdout"`
	Debug     bool   `yaml:"debug"`
}

// PublicAPIServerConfig contains configuration to provide public REST API.
type PublicAPIServerConfig struct {
	ServerAddress string `yaml:"server_address"`
	ServerPort    int    `yaml:"server_port"`
	ReadTimeout   int    `yaml:"read_timeout"`
	WriteTimeout  int    `yaml:"write_timeout"`
	IdleTimeout   int    `yaml:"idle_timeout"`
}

// ServiceAPIServerConfig contains configuration to provide service REST API.
type ServiceAPIServerConfig struct {
	ServerAddress string `yaml:"server_address"`
	ServerPort    int    `yaml:"server_port"`
	ReadTimeout   int    `yaml:"read_timeout"`
	WriteTimeout  int    `yaml:"write_timeout"`
	IdleTimeout   int    `yaml:"idle_timeout"`
}

// CheckConfig helps to check if global application config is ready.
func CheckConfig() error {
	if Config == nil {
		return errors.New(errGlobalConfig)
	}

	return nil
}

// InitFromFile reads config file and initializes global config.
func InitFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := initFromString(data); err != nil {
		return err
	}

	log.Printf("Config loaded from: %s", path)

	return nil
}

// initFromString reads raw string and initializes global config.
func initFromString(data []byte) error {
	cfg := AppConfig{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	Config = &cfg

	// Set default string parameters if omitted.
	defaultStringParameters := map[*string]string{
		&Config.PublicAPI.ServerAddress:  defaultPublicAPIAddress,
		&Config.ServiceAPI.ServerAddress: defaultServiceAPIAddress,
	}
	for currentValue, defaultValue := range defaultStringParameters {
		setDefaultStringValue(currentValue, defaultValue)
	}

	// Set default int parameters if omitted.
	defaultIntParameters := map[*int]int{
		// Public API defaults
		&Config.PublicAPI.ServerPort:   defaultPublicAPIPort,
		&Config.PublicAPI.ReadTimeout:  defaultHTTPReadTimeout,
		&Config.PublicAPI.WriteTimeout: defaultHTTPWriteTimeout,
		&Config.PublicAPI.IdleTimeout:  defaultHTTPIdleTimeout,
		// ServiceAPI defaults
		&Config.ServiceAPI.ServerPort:   defaultServiceAPIPort,
		&Config.ServiceAPI.ReadTimeout:  defaultHTTPReadTimeout,
		&Config.ServiceAPI.WriteTimeout: defaultHTTPWriteTimeout,
		&Config.ServiceAPI.IdleTimeout:  defaultHTTPIdleTimeout,
	}
	for currentValue, defaultValue := range defaultIntParameters {
		setDefaultIntValue(currentValue, defaultValue)
	}

	return nil
}

func setDefaultIntValue(currentValue *int, defaultValue int) {
	if *currentValue <= 0 {
		*currentValue = defaultValue
	}
}

func setDefaultStringValue(currentValue *string, defaultValue string) {
	if *currentValue == "" {
		*currentValue = defaultValue
	}
}
