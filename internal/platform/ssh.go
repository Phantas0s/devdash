package platform

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	sshAgentEnv = "SSH_AUTH_SOCK"
)

type SSH struct {
	client *ssh.Client
}

// TODO accept more method of connection than only via ssh-agent
func NewSSHClient(username, addr string) (*SSH, error) {
	sshClient, err := createAgentAuth(username, addr)
	if err != nil {
		return nil, err
	}
	return &SSH{
		client: sshClient,
	}, nil
}

func createAgentAuth(username, addr string) (*ssh.Client, error) {
	s := os.Getenv(sshAgentEnv)
	if s == "" {
		return nil, errors.Errorf("%s environment varible empty", sshAgentEnv)
	}

	agentConn, err := net.Dial("unix", s)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't connect via ssh-agent")
	}

	auth := ssh.PublicKeysCallback(agent.NewClient(agentConn).Signers)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{auth},
	}

	return ssh.Dial("tcp", addr, config)
}

// Run a command on remote server via SSH
func (s *SSH) Run(command string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", errors.Wrapf(err, "can't create session with SSH client for command %s", command)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		return "", errors.Wrapf(err, "can't run command %s on remote server", command)
	}

	return string(buf.Bytes()), nil
}

func (s *SSH) getMemoryInfo(headers []string, metrics []string) (cells [][]string, err error) {
	lines, err := s.Run("/bin/cat /proc/meminfo")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	var data string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			val = incSizeMetric(val)
			if err != nil {
				data += "unknown" + ","
				continue
			}

			for _, v := range metrics {
				if parts[0] == v {
					data += strconv.FormatInt(int64(val), 10) + ","
				}
			}
		}
	}

	cells = append(cells, headers)
	for _, v := range formatToTable(headers, data) {
		cells = append(cells, v)
	}

	return cells, nil
}

// formatToTable display.
// The string needs to have this:
// Info needs to be splitted with comma
// Depending on number of columns (headers)
// TODO improve this comment :D
func formatToTable(headers []string, data string) (cells [][]string) {
	col := len(headers)
	c := strings.Split(data, ",")
	lenCells := len(c)
	for i := 0; i < lenCells; i += col {
		cells = append(cells, c[i:min(i+col, lenCells)])
	}

	return cells
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func incSizeMetric(val uint64) uint64 {
	return val / 1024
}
