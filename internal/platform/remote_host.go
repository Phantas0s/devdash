package platform

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/gokit"
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
// TODO should be possible to put run as a property of remote host object
// Could swap between localhost and remotehost that way when call NewRemoteHost with a flag
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
	secs, err = strconv.ParseFloat(d[0], strconv.IntSize)
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
// TODO not used for now - create a table from the return of correctly formatted any command
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

		val, err := strconv.ParseUint(parts[1], 10, strconv.IntSize)
		if err != nil {
			return nil, err
		}

		for _, v := range metrics {
			if strings.Trim(parts[0], ":") == v {
				data += strconv.FormatUint(val, 10) + ","
			}
		}
	}

	values := formatToBar(data)
	result := []int{}
	for _, v := range values {
		result = append(result, int(gokit.ConvertBinUnit(float64(v), "kb", unit)))
	}

	return result, nil
}

func (s *RemoteHost) MemoryRate() (float64, error) {
	lines, err := s.run("/bin/cat /proc/meminfo")
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	var memTotal float64 = 0
	var memFree float64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 3 {
			continue
		}

		val, err := strconv.ParseFloat(parts[1], strconv.IntSize)
		if err != nil {
			return 0, err
		}

		k := strings.Trim(parts[0], ":")
		switch k {
		case "MemTotal":
			memTotal = gokit.ConvertBinUnit(val, "kb", "mb")
		case "MemFree":
			memFree = gokit.ConvertBinUnit(val, "kb", "mb")
		}

	}

	memUsed := memTotal - memFree
	return gokit.Round(float64(memUsed)*100/float64(memTotal), 2), nil
}

// TODO to refactor - DRY
func (s *RemoteHost) SwapRate() (float64, error) {
	lines, err := s.run("/bin/cat /proc/meminfo")
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	var swapTotal float64 = 0
	var swapFree float64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 3 {
			continue
		}

		val, err := strconv.ParseFloat(parts[1], strconv.IntSize)
		if err != nil {
			return 0, err
		}

		k := strings.Trim(parts[0], ":")
		switch k {
		case "SwapTotal":
			swapTotal = gokit.ConvertBinUnit(val, "kb", "mb")
		case "SwapFree":
			swapFree = gokit.ConvertBinUnit(val, "kb", "mb")
		}

	}

	swapUsed := swapTotal - swapFree
	return gokit.Round(float64(swapUsed)*100/float64(swapTotal), 2), nil
}

// See https://www.idnt.net/en-US/kb/941772
func (s *RemoteHost) CPURate() (float64, error) {
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
	return 100 - gokit.Round(float64(idle)*100/float64(total), 2), nil
}

// GetNetStat returns net stat
func (s *RemoteHost) NetIO(unit string) (string, error) {
	lines, err := s.run("/bin/cat /proc/net/dev")
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	var receiveBytes uint64 = 0
	var transmitBytes uint64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 10 {
			continue
		}

		device := strings.TrimSpace(strings.Trim(parts[0], ":"))

		if device != "lo" {
			rb, _ := strconv.ParseUint(parts[2], 10, strconv.IntSize)
			tb, _ := strconv.ParseUint(parts[10], 10, strconv.IntSize)

			receiveBytes += rb
			transmitBytes += tb
		}
	}

	rx := strconv.FormatFloat(gokit.ConvertBinUnit(float64(receiveBytes), "b", unit), 'f', 2, strconv.IntSize)
	tx := strconv.FormatFloat(gokit.ConvertBinUnit(float64(transmitBytes), "b", unit), 'f', 2, strconv.IntSize)

	return rx + " / " + tx, nil
}

func (s *RemoteHost) Disk(headers []string, unit string) ([][]string, error) {
	// GetIOStat returns io stat
	lines, err := s.run("/bin/df -x devtmpfs -x tmpfs -x debugfs")
	if err != nil {
		return nil, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	data := ""
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 6 || count == 0 {
			count++
			continue
		}

		filesystem := parts[0]
		size, _ := strconv.ParseFloat(parts[1], strconv.IntSize)
		used, _ := strconv.ParseFloat(parts[2], strconv.IntSize)
		available, _ := strconv.ParseFloat(parts[3], strconv.IntSize)
		useRate := parts[4]
		mount := parts[5]

		d := []string{
			filesystem,
			strconv.FormatFloat(gokit.ConvertBinUnit(size, "kb", unit), 'f', 2, strconv.IntSize) + unit,
			strconv.FormatFloat(gokit.ConvertBinUnit(used, "kb", unit), 'f', 2, strconv.IntSize) + unit,
			strconv.FormatFloat(gokit.ConvertBinUnit(available, "kb", unit), 'f', 2, strconv.IntSize) + unit,
			useRate,
			mount,
		}

		data += strings.Join(d, ",") + ","
		count++
	}

	c := [][]string{headers}
	c = append(c, formatToTable(headers, data)...)

	return c, nil
}

func (s *RemoteHost) DiskIO(unit string) (string, error) {
	// GetIOStat returns io stat
	lines, err := s.run("/bin/cat /proc/diskstats")
	if err != nil {
		return "", nil
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	var read uint64 = 0
	var write uint64 = 0

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 9 {
			continue
		}

		if parts[3] == "0" {
			continue
		}

		r, _ := strconv.ParseUint(parts[5], 10, strconv.IntSize)
		w, _ := strconv.ParseUint(parts[9], 10, strconv.IntSize)

		read += r * 512
		write += w * 512
	}

	fr := gokit.ConvertBinUnit(float64(read), "kb", unit)
	fw := gokit.ConvertBinUnit(float64(write), "kb", unit)

	return strconv.FormatFloat(fr, 'f', 2, strconv.IntSize) + " / " + strconv.FormatFloat(fw, 'f', 2, strconv.IntSize), nil
}

func formatToBar(data string) (val []uint64) {
	data = strings.Trim(data, ",")
	s := strings.Split(data, ",")
	val = []uint64{}
	for _, v := range s {
		k, _ := strconv.ParseUint(v, 10, strconv.IntSize)
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
		next := c[i:min(i+col, lenCells)]
		if len(next) == col {
			cells = append(cells, next)
		}
	}

	return cells
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
