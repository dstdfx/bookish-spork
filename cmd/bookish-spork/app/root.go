package app

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dstdfx/bookish-spork/internal/app/bookishspork"
	"github.com/dstdfx/bookish-spork/internal/pkg/config"
	"github.com/dstdfx/bookish-spork/internal/pkg/log"
	"github.com/spf13/cobra"
)

const defaultCfgFile = "/etc/bookish-spork/bookish-spork.yaml"

var cfgFile string

// Variables that are injected in build time.
var (
	buildGitCommit string
	buildGitTag    string
	buildDate      string
	buildCompiler  = runtime.Version()
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "bookish-spork",
	Short: "bookish-spork represents a simple HTTP API interface to in-memory cache",
	Run: func(_ *cobra.Command, _ []string) {
		// Initialize application config and log
		if _, err := os.Stat(cfgFile); err != nil {
			exitWithErr(fmt.Errorf("config file %s can't be read: %s", cfgFile, err))
		}
		if err := config.InitFromFile(cfgFile); err != nil {
			exitWithErr(err)
		}

		// Init logger
		logger, err := log.InitLogger(log.InitLoggerOpts{
			File:      config.Config.Log.File,
			UseStdout: config.Config.Log.UseStdout,
			Debug:     config.Config.Log.Debug,
		})
		if err != nil {
			exitWithErr(err)
		}

		opts := bookishspork.StartOpts{
			Interrupt:      make(chan os.Signal, 1),
			BuildGitCommit: buildGitCommit,
			BuildGitTag:    buildGitTag,
			BuildDate:      buildDate,
			BuildCompiler:  buildCompiler,
		}

		//Start main routine
		if err := bookishspork.StartService(logger, opts); err != nil {
			exitWithErr(fmt.Errorf("error starting bookish-spork app: %w", err))
		}
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config",
		defaultCfgFile, "path to application config")
}

// exitWithErr is a helper method to print errors in case of empty logger.
func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "application is exiting after error: %s\n", err)
	os.Exit(1)
}
