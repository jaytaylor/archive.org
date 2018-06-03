package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"jaytaylor.com/archive.org"
)

var (
	Quiet          bool
	Verbose        bool
	Wait           bool
	RequestTimeout time.Duration = archiveorg.DefaultRequestTimeout
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "Activate quiet log output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose log output")
	rootCmd.PersistentFlags().DurationVarP(&RequestTimeout, "request-timeout", "r", RequestTimeout, "Timeout duration for HTTP requests")
	rootCmd.PersistentFlags().StringVarP(&archiveorg.BaseURL, "base-url", "b", archiveorg.BaseURL, "Archive.org server base URL address")
	rootCmd.PersistentFlags().StringVarP(&archiveorg.HTTPHost, "http-host", "", archiveorg.HTTPHost, "'Host' header to use")
	rootCmd.PersistentFlags().StringVarP(&archiveorg.UserAgent, "user-agent", "u", archiveorg.UserAgent, "'User-Agent' header to use")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		errorExit(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "archive.org",
	Short: "create a new snapshot of a URL",
	Long:  "command-line interface for requesting archive.org perform a fresh crawl of a URL",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(_ *cobra.Command, _ []string) {
		initLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		location, err := archiveorg.Capture(args[0], RequestTimeout)
		if err != nil {
			errorExit(err)
		}
		fmt.Println(location)
	},
}

func errorExit(err interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

func initLogging() {
	level := log.InfoLevel
	if Verbose {
		level = log.DebugLevel
	}
	if Quiet {
		level = log.ErrorLevel
	}
	log.SetLevel(level)
}
