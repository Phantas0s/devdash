package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List your devdash boards",
		// TODO write some help
		Run: func(cmd *cobra.Command, args []string) {
			runList()
		},
	}

	return listCmd
}

// TODO to refactor and TEST
func runList() {
	homeFiles, err := ioutil.ReadDir(dashPath())
	if err != nil {
		log.Fatal(err)
	}
	currentFiles, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	fs := []fs.FileInfo{}
	for _, v := range homeFiles {
		f, isDash := isDashboard(v, dashPath())
		if isDash {
			fs = append(fs, f)
		}
	}

	for _, v := range currentFiles {
		f, isDash := isDashboard(v, ".")
		if isDash {
			fs = append(fs, f)
		}
	}

	for _, f := range fs {
		s := strings.Split(f.Name(), ".")
		// TODO erk
		if !f.IsDir() && len(s) > 1 && (s[1] == "json" || s[1] == "toml" || s[1] == "yaml" || s[1] == "yml") {
			fmt.Fprintln(os.Stdout, s[0])
		}
	}
}

func isDashboard(fileInfo fs.FileInfo, path string) (fs.FileInfo, bool) {
	file, err := os.Open(path + string(filepath.Separator) + fileInfo.Name())
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "projects:") {
			return fileInfo, true
		}
	}

	return nil, false
}
