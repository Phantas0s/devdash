package platform

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

const (
	sshAgentEnv = "SSH_AUTH_SOCK"
)

type RemoteHost struct {
	sshClient *ssh.Client
}

// TODO accept more method of connection than only via ssh-agent
func NewRemoteHost(username, addr string) (*RemoteHost, error) {
	sshClient, err := sshAgentAuth(username, addr)
	if err != nil {
		return nil, err
	}
	return &RemoteHost{
		sshClient: sshClient,
	}, nil
}

// Run a command on remote server via SSH
func (s *RemoteHost) run(command string) (string, error) {
	session, err := s.sshClient.NewSession()
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

func (s *RemoteHost) Uptime() (int64, error) {
	command := "/bin/cat /proc/uptime"
	uptime, err := s.run(command)
	if err != nil {
		return 0, err
	}

	d := strings.Fields(uptime)
	if len(d) == 0 {
		return 0, errors.Errorf("command %s return nothing", command)
	}

	var secs float64
	secs, err = strconv.ParseFloat(d[0], 64)
	if err != nil {
		return 0, err
	}

	return int64(time.Duration(secs * 1e9)), nil
}

// TODO Doesn't really make sense to have memory info with headers on top
// it's more headers on the left and one value on the right
// Simple textbox for each of them would be enough?
// Also, stack bar for total memory / memory free | total swap / swap free could be nice
func (s *RemoteHost) getMemoryInfo(headers []string, metrics []string) (cells [][]string, err error) {
	lines, err := s.run("/bin/cat /proc/meminfo")
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

// TODO no need to specify headers
func (s *RemoteHost) Table(command string, headers []string, metrics []string) (cells [][]string, err error) {
	lines, err := s.run(command)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	data := ""
	for scanner.Scan() {
		line := scanner.Text()
		data += line
	}

	return formatToTable(headers, data), nil
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
