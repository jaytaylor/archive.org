package main

import (
	"encoding/json"
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
	Use:   "archive.org-snapshots",
	Short: "search for archive.org snapshots",
	Long:  "command-line interface for searching archive.org for URL page snapshots",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(_ *cobra.Command, _ []string) {
		initLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		snapshots, err := archiveorg.Search(args[0], RequestTimeout)
		if err != nil {
			errorExit(err)
		}

		log.Infof("Found %v results", len(snapshots))

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")

		if err := enc.Encode(&snapshots); err != nil {
			errorExit(fmt.Errorf("marshalling snapshots to JSON: %s", err))
		}
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
