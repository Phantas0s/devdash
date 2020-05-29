package gokit

import (
	"bytes"
	"os/exec"
	"strings"
)

func Pipeline(command string) (out, errs []byte, pipeLineError error) {
	cmds := []*exec.Cmd{}
	piped := strings.Split(command, "|")
	for _, v := range piped {
		c := strings.Split(strings.TrimSpace(v), " ")
		for k, v := range c[1:] {
			c[k+1] = strings.Replace(v, "'", "", -1)
		}
		// let it here for debug purposes... might need it :D
		// fmt.Println(c[0])
		// fmt.Println(strings.Join(c[1:], ","))
		cmds = append(cmds, exec.Command(strings.TrimSpace(c[0]), c[1:]...))
	}

	var stderr bytes.Buffer
	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		cmd.Stderr = &stderr
	}

	var output bytes.Buffer
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	return output.Bytes(), stderr.Bytes(), nil
}
