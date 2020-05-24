package platform

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/humanmath"
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

func (s *RemoteHost) Load() (string, error) {
	command := "/bin/cat /proc/loadavg"
	lines, err := s.run(command)
	if err != nil {
		return "", err
	}
	res := strings.Fields(lines)
	if len(res) < 3 {
		return "", errors.Errorf(
			"command %s return unexpected %v, needs to have 3 parts separated with whitespaces",
			command,
			res,
		)
	}

	return fmt.Sprintf("%s %s %s", res[0], res[1], res[2]), nil
}

func (s *RemoteHost) Processes() (string, error) {
	command := "/bin/cat /proc/loadavg"
	lines, err := s.run(command)
	if err != nil {
		return "", err
	}

	res := strings.Fields(lines)
	if len(res) < 5 {
		return "", errors.Errorf("command %s return unexpected %v, needs to have 5 parts separated with whitespaces", command, res)
	}

	runProc := "unknown"
	totalProc := "unknown"
	if i := strings.Index(res[3], "/"); i != -1 {
		runProc = res[3][0:i]
		if i+1 < len(res[3]) {
			totalProc = res[3][i+1:]
		}
	}

	return fmt.Sprintf("%s/%s", runProc, totalProc), nil
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

	return humanmath.ConvertBinUnit(formatToBar(data), "kb", unit), nil
}

// See https://www.idnt.net/en-US/kb/941772
func (s *RemoteHost) CPURate() (int, error) {
	raw, err := s.run("/bin/cat /proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(raw), "\n")

	// aggregate of all other cpus
	cpu := strings.Fields(lines[0])

	if len(cpu) < 5 {
		return 0, errors.Errorf("needs 5 fields for cpu: header, user, nice, system, idle. Instead, having %s", cpu)
	}

	user, _ := strconv.ParseUint(cpu[1], 10, strconv.IntSize)
	nice, _ := strconv.ParseUint(cpu[2], 10, strconv.IntSize)
	system, _ := strconv.ParseUint(cpu[3], 10, strconv.IntSize)
	idle, _ := strconv.ParseUint(cpu[4], 10, strconv.IntSize)

	var IOWait uint64 = 0
	var IRQ uint64 = 0
	var softIRQs uint64 = 0
	var steal uint64 = 0
	var guest uint64 = 0
	var guestNice uint64 = 0
	if len(cpu) > 5 {
		IOWait, _ = strconv.ParseUint(cpu[5], 10, strconv.IntSize)
	}
	if len(cpu) > 6 {
		IRQ, _ = strconv.ParseUint(cpu[6], 10, strconv.IntSize)
	}
	if len(cpu) > 7 {
		softIRQs, _ = strconv.ParseUint(cpu[7], 10, strconv.IntSize)
	}
	if len(cpu) > 8 {
		steal, _ = strconv.ParseUint(cpu[8], 10, strconv.IntSize)
	}
	if len(cpu) > 9 {
		guest, _ = strconv.ParseUint(cpu[9], 10, strconv.IntSize)
	}
	if len(cpu) > 10 {
		guestNice, _ = strconv.ParseUint(cpu[10], 10, strconv.IntSize)
	}

	total := user + nice + system + idle + IOWait + IRQ + softIRQs + steal + guest + guestNice

	// Percentage of not idle (busy)
	return int(100 - humanmath.Round(float64(idle)*100/float64(total), 2)), nil
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
