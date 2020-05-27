package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags
	cfg     string
	logpath string
	debug   bool

	rootCmd = &cobra.Command{
		Use:   "devdash",
		Short: "DevDash is a highly configurable terminal dashboard for developers",
		Long:  `DevDash is a highly flexible terminal dashboard for developers, which allows you to gather and refresh the data you really need from Google Analytics, Google Search Console, Github, TravisCI, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg, "config", "c", "", "A valid dashboard configuration")
	// rootCmd.PersistentFlags().StringVarP(&logpath, "logpath", "l", "", "Path for logging")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug Mode - doesn't display graph")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}

func run() {
	begin(cfg, debug)
}
