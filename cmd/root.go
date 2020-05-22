package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

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

// TODO - Init logger somewhere. Logger file tmp by default?
func InitLoggerFile(logpath string) *log.Logger {
	if logpath == "" {
		return log.New(os.Stderr, "", 0)
	}

	file, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(file, "", 0)
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " - ")

	return l
}
