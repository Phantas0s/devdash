package platform

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const git = "git"

type Git struct {
	Path string
}

func NewGit(path string) *Git {
	return &Git{
		Path: path,
	}
}

func (g *Git) Branches() ([][]string, error) {
	cmd := exec.Command(
		git,
		"for-each-ref",
		"--sort=committerdate",
		"refs/heads/",
		"--format=%(refname:short),%(authorname),%(authoremail),%(authordate:short)",
	)
	cmd.Dir = g.Path

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "can't run %v", git+" "+strings.Join(cmd.Args, " "))
	}

	output := cmdOutput.Bytes()
	return formatBranches(string(output)), nil
}

func formatBranches(data string) [][]string {
	data = strings.ReplaceAll(data, "<", "")
	data = strings.ReplaceAll(data, ">", "")
	result := [][]string{{"Branch", "Last Commit By", "Creator email", "Creation date"}}
	d := strings.Split(data, "\n")
	for i := 0; i < len(d)-1; i++ {
		result = append(result, strings.Split(d[i], ","))
	}

	fmt.Println(result)

	return result
}
