package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/pkg/errors"
)

const (
	rhUptime        = "rh.box_uptime"
	rhLoad          = "rh.box_load"
	rhProcesses     = "rh.box_processes"
	rhBoxMemRate    = "rh.box_memory_rate"
	rhGaugeMemRate  = "rh.gauge_memory_rate"
	rhBoxSwapRate   = "rh.box_swap_rate"
	rhGaugeSwapRate = "rh.gauge_swap_rate"
	rhBoxNetIO      = "rh.box_net_io"
	rhBoxDiskIO     = "rh.box_disk_io"
	rhBoxCPURate    = "rh.box_cpu_rate"
	rhGaugeCPURate  = "rh.gauge_cpu_rate"
	rhBarMemory     = "rh.bar_memory"
	rhBarRates      = "rh.bar_rates"
	rhTableDisk     = "rh.table_disk"
	rhTable         = "rh.table"
	rhBox           = "rh.box"
	rhGauge         = "rh.gauge"
	rhBar           = "rh.bar"
)

type HostWidget struct {
	tui     *Tui
	service *platform.Host
}

func NewHostWidget(username, addr string) (*HostWidget, error) {
	service, err := platform.NewHost(username, addr)
	if err != nil {
		return nil, err
	}

	return &HostWidget{
		service: service,
	}, nil
}

func (ms *HostWidget) CreateWidgets(widget Widget, tui *Tui) (f func() error, err error) {
	ms.tui = tui

	// Compatibility with localhost
	name := strings.Replace(widget.Name, "lh", "rh", 1)

	switch name {
	case rhUptime:
		f, err = ms.boxUptime(widget)
	case rhLoad:
		f, err = ms.boxLoad(widget)
	case rhProcesses:
		f, err = ms.boxProcesses(widget)
	case rhBarMemory:
		f, err = ms.barMemory(widget)
	case rhBoxCPURate:
		f, err = ms.boxCPURate(widget)
	case rhGaugeCPURate:
		f, err = ms.gaugeCPURate(widget)
	case rhBoxMemRate:
		f, err = ms.boxMemRate(widget)
	case rhGaugeMemRate:
		f, err = ms.gaugeMemRate(widget)
	case rhBoxSwapRate:
		f, err = ms.boxSwapRate(widget)
	case rhGaugeSwapRate:
		f, err = ms.gaugeSwapRate(widget)
	case rhBoxNetIO:
		f, err = ms.boxNetIO(widget)
	case rhBoxDiskIO:
		f, err = ms.boxDiskIO(widget)
	case rhBarRates:
		f, err = ms.barRates(widget)
	case rhTableDisk:
		f, err = ms.tableDisk(widget)
	case rhTable:
		f, err = ms.table(widget)
	case rhBox:
		f, err = ms.box(widget)
	case rhGauge:
		f, err = ms.gauge(widget)
	case rhBar:
		f, err = ms.bar(widget)
	default:
		return nil, errors.Errorf("can't find the widget %s", widget.Name)
	}
	return
}

func (ms *HostWidget) boxLoad(widget Widget) (f func() error, err error) {
	title := " Load "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	load, err := platform.HostLoad(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(load, title, widget.Options)
	}

	return
}

func (ms *HostWidget) boxProcesses(widget Widget) (f func() error, err error) {
	title := " Running processes "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	procs, err := platform.HostProcesses(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(procs, title, widget.Options)
	}

	return
}

func (ms *HostWidget) boxUptime(widget Widget) (f func() error, err error) {
	title := " Uptime "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	uptime, err := platform.HostUptime(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(formatSeconds(time.Duration(uptime)), title, widget.Options)
	}

	return
}

func formatSeconds(dur time.Duration) string {
	dur = dur - (dur % time.Second)
	var days int
	for dur.Hours() > 24.0 {
		days++
		dur -= 24 * time.Hour
	}
	for dur.Hours() > 24.0 {
		days++
		dur -= 24 * time.Hour
	}

	s1 := dur.String()
	s2 := ""
	if days > 0 {
		s2 = fmt.Sprintf("%dd ", days)
	}
	for _, ch := range s1 {
		s2 += string(ch)
		if ch == 'h' || ch == 'm' {
			s2 += " "
		}
	}
	return s2
}

func (ms *HostWidget) boxCPURate(widget Widget) (f func() error, err error) {
	title := " CPU usage "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	CPURate, err := platform.HostCPURate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(strconv.FormatFloat(CPURate, 'f', 2, strconv.IntSize)+" %", title, widget.Options)
	}

	return
}

func (ms *HostWidget) gaugeCPURate(widget Widget) (f func() error, err error) {
	title := " CPU usage "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	CPURate, err := platform.HostCPURate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddGauge(CPURate, title, widget.Options)
	}

	return
}

func (ms *HostWidget) boxMemRate(widget Widget) (f func() error, err error) {
	title := " Memory usage "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	memRate, err := platform.HostMemoryRate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(strconv.FormatFloat(memRate, 'f', 2, 64)+" %", title, widget.Options)
	}

	return
}

func (ms *HostWidget) gaugeMemRate(widget Widget) (f func() error, err error) {
	title := " Memory usage "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	memRate, err := platform.HostMemoryRate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddGauge(memRate, title, widget.Options)
	}

	return
}

func (ms *HostWidget) boxSwapRate(widget Widget) (f func() error, err error) {
	title := " Swap usage "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	swapRate, err := platform.HostSwapRate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(strconv.FormatFloat(swapRate, 'f', 2, 64)+" %", title, widget.Options)
	}

	return
}

func (ms *HostWidget) gaugeSwapRate(widget Widget) (f func() error, err error) {
	title := " Swap usage "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	swapRate, err := platform.HostSwapRate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddGauge(swapRate, title, widget.Options)
	}

	return
}

func (ms *HostWidget) barRates(widget Widget) (f func() error, err error) {
	title := " Resources usage (%) "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	swapRate, err := platform.HostSwapRate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	cpuRate, err := platform.HostCPURate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	memoryRate, err := platform.HostMemoryRate(ms.service.Runner)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddBarChart(
			[]int{int(cpuRate), int(memoryRate), int(swapRate)},
			[]string{"CPU", "Memory", "Swap"},
			title,
			widget.Options,
		)
	}

	return
}

func (ms *HostWidget) boxNetIO(widget Widget) (f func() error, err error) {
	unit := "kb"
	if _, ok := widget.Options[optionUnit]; ok {
		unit = widget.Options[optionUnit]
	}

	title := fmt.Sprintf(" Net I/O (%s) ", strings.ToUpper(unit))
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	netIO, err := platform.HostNetIO(ms.service.Runner, unit)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(netIO, title, widget.Options)
	}

	return
}

func (ms *HostWidget) boxDiskIO(widget Widget) (f func() error, err error) {
	unit := "kb"
	if _, ok := widget.Options[optionUnit]; ok {
		unit = widget.Options[optionUnit]
	}

	title := fmt.Sprintf(" Disk I/O (%s) ", strings.ToUpper(unit))
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	diskIO, err := platform.HostDiskIO(ms.service.Runner, unit)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(diskIO, title, widget.Options)
	}

	return
}

func (ms *HostWidget) barMemory(widget Widget) (f func() error, err error) {
	metrics := []string{"MemTotal", "MemFree", "MemAvailable"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(strings.TrimSpace(widget.Options[optionMetrics]), ",")
		}
	}

	headers := metrics
	if _, ok := widget.Options[optionHeaders]; ok {
		if len(widget.Options[optionHeaders]) > 0 {
			headers = strings.Split(strings.TrimSpace(widget.Options[optionHeaders]), ",")
		}
	}

	unit := "kb"
	if _, ok := widget.Options[optionUnit]; ok {
		unit = widget.Options[optionUnit]
	}

	title := fmt.Sprintf(" Memory (%s) ", strings.ToUpper(unit))
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	mem, err := platform.HostMemory(ms.service.Runner, metrics, unit)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddBarChart(mem, headers, title, widget.Options)
	}

	return
}

func (ms *HostWidget) tableDisk(widget Widget) (f func() error, err error) {
	unit := "gb"
	if _, ok := widget.Options[optionUnit]; ok {
		unit = widget.Options[optionUnit]
	}

	title := " Disks "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	headers := []string{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"}
	if _, ok := widget.Options[optionHeaders]; ok {
		if len(widget.Options[optionHeaders]) > 0 {
			headers = strings.Split(strings.TrimSpace(widget.Options[optionHeaders]), ",")
		}
	}

	data, err := platform.HostDisk(ms.service.Runner, headers, unit)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTable(data, title, widget.Options)
	}

	return
}

func (ms *HostWidget) table(widget Widget) (f func() error, err error) {
	title := " Table "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	headers := []string{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"}

	cmd := "/bin/df -x devtmpfs -x tmpfs -x debugfs | tail -n +2"
	if _, ok := widget.Options[optionCommand]; ok {
		cmd = widget.Options[optionCommand]
		headers = []string{}
	}

	if _, ok := widget.Options[optionHeaders]; ok {
		if len(widget.Options[optionHeaders]) > 0 {
			headers = strings.Split(strings.TrimSpace(widget.Options[optionHeaders]), ",")
		}
	}

	data, err := platform.HostTable(ms.service.Runner, cmd, headers)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTable(data, title, widget.Options)
	}

	return
}

func (ms *HostWidget) box(widget Widget) (f func() error, err error) {
	title := " Box "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	cmd := "echo 'box'"
	if _, ok := widget.Options[optionCommand]; ok {
		cmd = widget.Options[optionCommand]
	}

	data, err := platform.HostBox(ms.service.Runner, cmd)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(data, title, widget.Options)
	}

	return
}

func (ms *HostWidget) gauge(widget Widget) (f func() error, err error) {
	title := " Gauge "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	cmd := "echo 50"
	// cmd := "/bin/df -x devtmpfs -x tmpfs -x debugfs | sed -n '1!p'"
	if _, ok := widget.Options[optionCommand]; ok {
		cmd = widget.Options[optionCommand]
	}

	data, err := platform.HostGauge(ms.service.Runner, cmd)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddGauge(data, title, widget.Options)
	}

	return
}

func (ms *HostWidget) bar(widget Widget) (f func() error, err error) {
	title := " Example of bar "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	headers := []string{"small", "bigger", "big", "insane"}
	if _, ok := widget.Options[optionHeaders]; ok {
		if len(widget.Options[optionHeaders]) > 0 {
			headers = strings.Split(strings.TrimSpace(widget.Options[optionHeaders]), ",")
		}
	}

	cmd := "echo -e 20 30 40 50"
	if _, ok := widget.Options[optionCommand]; ok {
		cmd = widget.Options[optionCommand]
	}

	data, err := platform.HostBar(ms.service.Runner, cmd)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddBarChart(
			data,
			headers,
			title,
			widget.Options,
		)
	}

	return
}
