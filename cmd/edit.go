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
		Short: "Edit your dashboard with an shell editor",
		Run: func(cmd *cobra.Command, args []string) {
			edit(args)
		},
	}

	editCmd.Flags().StringVarP(&editor, "editor", "e", "$EDITOR", "Path of your favorite editor")

	return editCmd
}

func edit(args []string) {
	file := findConfigFile(args[0])
	if file == "" {
		fmt.Fprintf(os.Stdout, "The config %s doesn't exist", args[0])
		return
	} else {
		editDashboard(os.ExpandEnv(editor), filepath.Join(dashPath(), file))
	}
}

// TODO add that to the gokit/cmd.
func editDashboard(editor string, config string) {
	cmd := exec.Command(editor, config)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func findConfigFile(search string) string {
	fs := getConfigFiles()
	for _, v := range fs {
		if search == removeExt(v.Name()) || search == v.Name() {
			return v.Name()
		}
	}

	return ""
}
