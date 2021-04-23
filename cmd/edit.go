package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var editor string

func editCmd() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "edit dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			edit(args)
		},
	}

	editCmd.Flags().StringVarP(&editor, "editor", "e", "$EDITOR", "Path of your favorite editor")

	return editCmd
}

func edit(args []string) {
	// TODO make it work for every OS
	cmd := exec.Command(os.ExpandEnv(editor), filepath.Join(dashPath(), args[0]))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
