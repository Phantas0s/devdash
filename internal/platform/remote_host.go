package platform

import (
	"bufio"
	"bytes"
	"math"
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
	if len(d) < 1 {
		return 0, errors.Errorf("command %s return nothing", command)
	}

	var secs float64
	secs, err = strconv.ParseFloat(d[0], 64)
	if err != nil {
		return 0, err
	}

	return int64(time.Duration(secs * 1e9)), nil
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

func (s *RemoteHost) Memory(metrics []string, unit string) (val []int, err error) {
	lines, err := s.run("/bin/cat /proc/meminfo")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	var data string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 3 {
			continue
		}

		val, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}

		for _, v := range metrics {
			if strings.Trim(parts[0], ":") == v {
				data += strconv.FormatInt(int64(val), 10) + ","
			}
		}
	}

	return ConvertBinUnit(formatToBar(data), "kb", unit), nil
}

// TODO these conversions are pretty ugly
// TODO floating point number (precision 2)
func ConvertBinUnit(val []int, base, to string) (tv []int) {
	convert := map[string]int{
		"b":  3,
		"kb": 2,
		"mb": 1,
		"gb": 0,
	}

	r := convert[base] - convert[to]

	if r > 0 {
		for _, v := range val {
			tv = append(tv, int(math.Floor(float64(v)/(math.Pow(float64(1024), float64(r))))))
		}
	}

	if r < 0 {
		for _, v := range val {
			tv = append(tv, int(math.Floor(float64(v)*(math.Pow(float64(1024), float64(r))))))
		}
	}

	return tv
}

func formatToBar(data string) (val []int) {
	data = strings.Trim(data, ",")
	s := strings.Split(data, ",")
	val = []int{}
	for _, v := range s {
		k, _ := strconv.Atoi(v)
		val = append(val, k)
	}

	return
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
