package platform

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func Test_formatToTable(t *testing.T) {
	testCases := []struct {
		name     string
		expected [][]string
		headers  []string
		data     string
	}{
		{
			name: "happy case",
			expected: [][]string{
				{"data1", "data2", "data3", "data4"},
				{"data5", "data6", "data7", "data8"},
			},
			headers: []string{"dataHeader1", "dataHeader2", "dataHeader3", "dataHeader4"},
			data:    "data1,data2,data3,data4,data5,data6,data7,data8",
		},
		{
			name: "data dropped when last row too small",
			expected: [][]string{
				{"data1", "data2", "data3", "data4"},
				{"data5", "data6", "data7", "data8"},
			},
			headers: []string{"dataHeader1", "dataHeader2", "dataHeader3", "dataHeader4"},
			data:    "data1,data2,data3,data4,data5,data6,data7,data8,data9,data10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatToTable(len(tc.headers), tc.data)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_formatToBar(t *testing.T) {
	testCases := []struct {
		name     string
		expected []uint64
		data     string
	}{
		{
			name: "happy case",
			expected: []uint64{
				10, 20, 30, 40,
			},
			data: "10,20,30,40",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatToBar(tc.data)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostUptime(t *testing.T) {
	testCases := []struct {
		name     string
		expected int64
		runner   runnerFunc
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: 17200210000000,
			runner:   func(cmd string) (string, error) { return "17200.21 59425.48", nil },
			wantErr:  false,
		},
		{
			name:    "Empty result",
			runner:  func(cmd string) (string, error) { return "", nil },
			wantErr: true,
		},
		{
			name:    "Runner return error",
			runner:  func(cmd string) (string, error) { return "", errors.New("Error!") },
			wantErr: true,
		},
		{
			name:    "Runner return wrong number",
			runner:  func(cmd string) (string, error) { return "hello", nil },
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostUptime(tc.runner)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostLoad(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		runner   runnerFunc
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: "1.12 1.21 0.96",
			runner:   func(cmd string) (string, error) { return "1.12 1.21 0.96 3/760 18313", nil },
			wantErr:  false,
		},
		{
			name:    "empty result",
			runner:  func(cmd string) (string, error) { return "", nil },
			wantErr: true,
		},
		{
			name:    "runner return error",
			runner:  func(cmd string) (string, error) { return "1.12 1.21 0.96 3/760 18313", errors.New("ERROR") },
			wantErr: true,
		},
		{
			name:     "runner return wrong result",
			expected: "",
			runner:   func(cmd string) (string, error) { return "hello", nil },
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostLoad(tc.runner)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
func Test_HostProcesses(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		runner   runnerFunc
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: "3/760",
			runner:   func(cmd string) (string, error) { return "1.12 1.21 0.96 3/760 18313", nil },
			wantErr:  false,
		},
		{
			name:    "empty result",
			runner:  func(cmd string) (string, error) { return "", nil },
			wantErr: true,
		},
		{
			name:    "runner return error",
			runner:  func(cmd string) (string, error) { return "1.12 1.21 0.96 3/760 18313", errors.New("ERROR") },
			wantErr: true,
		},
		{
			name:    "runner return wrong result",
			runner:  func(cmd string) (string, error) { return "hello", nil },
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostProcesses(tc.runner)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostMemory(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected []int
		metrics  []string
		unit     string
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: []int{8037936, 1423776, 3701620},
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_memory", t)), nil
			},
			metrics: []string{"MemTotal", "MemFree", "MemAvailable"},
			unit:    "kb",
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: []int{0},
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			metrics: []string{"MemTotal", "MemFree", "MemAvailable"},
			unit:    "kb",
			wantErr: false,
		},
		{
			name: "runner return error",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_memory", t)), errors.New("Error")
			},
			metrics: []string{"MemTotal", "MemFree", "MemAvailable"},
			unit:    "kb",
			wantErr: true,
		},
		{
			name: "runner return wrong result",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/ga_users.json", t)), errors.New("Error")
			},
			metrics: []string{"MemTotal", "MemFree", "MemAvailable"},
			unit:    "kb",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostMemory(tc.runner, tc.metrics, tc.unit)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostMemoryRate(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected float64
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: 82.29,
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_memory", t)), nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: 0,
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name: "runner return error",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_memory", t)), errors.New("Error")
			},
			wantErr: true,
		},
		{
			name: "runner return wrong result",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/ga_users.json", t)), errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostMemoryRate(tc.runner)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostSwapRate(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected float64
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: 5.96,
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_memory", t)), nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: 0,
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name: "runner return error",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_memory", t)), errors.New("Error")
			},
			wantErr: true,
		},
		{
			name: "runner return wrong result",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/ga_users.json", t)), errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostSwapRate(tc.runner)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostCPURate(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected float64
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: 11.969999999999999,
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_cpu", t)), nil
			},
			wantErr: false,
		},
		{
			name: "empty result",
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: true,
		},
		{
			name: "runner return error",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_cpu", t)), errors.New("Error")
			},
			wantErr: true,
		},
		{
			name: "runner return wrong result",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/ga_users.json", t)), errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostCPURate(tc.runner)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostNetIO(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected string
		unit     string
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: "322.29 / 146.11",
			unit:     "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_net", t)), nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: "0.00 / 0.00",
			unit:     "kb",
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name: "runner return error",
			unit: "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_net", t)), errors.New("Error")
			},
			wantErr: true,
		},
		{
			name: "runner return wrong result",
			unit: "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/ga_users.json", t)), errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostNetIO(tc.runner, tc.unit)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostDisk(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		headers  []string
		expected [][]string
		unit     string
		wantErr  bool
	}{
		{
			name:    "happy case",
			headers: []string{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"},
			expected: [][]string{
				{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"},
				{"/dev/sda3", "61665068.00kb", "16191716.00kb", "42311240.00kb", "28%", "/"},
				{"/dev/sda1", "194235.00kb", "54856.00kb", "125043.00kb", "31%", "/boot"},
			},
			unit: "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_disk", t)), nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			headers:  []string{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"},
			unit:     "kb",
			expected: [][]string{{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"}},
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name:    "runner return error",
			headers: []string{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"},
			unit:    "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_disk", t)), errors.New("Error")
			},
			wantErr: true,
		},
		{
			name:     "runner return wrong result",
			headers:  []string{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"},
			expected: [][]string{{"Filesystem", "Size", "Used", "Available", "Use%", "Mount"}},
			unit:     "kb",
			runner: func(cmd string) (string, error) {
				return "hello", nil
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostDisk(tc.runner, tc.headers, tc.unit)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostDiskIO(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected string
		unit     string
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: "4170343424.00 / 5843349504.00",
			unit:     "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_disk_io", t)), nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: "0.00 / 0.00",
			unit:     "kb",
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name: "runner return error",
			unit: "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_disk_io", t)), errors.New("Error")
			},
			wantErr: true,
		},
		{
			name:     "runner return wrong result",
			unit:     "kb",
			expected: "0.00 / 0.00",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/ga_users.json", t)), nil
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostDiskIO(tc.runner, tc.unit)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostBox(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected string
		command  string
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: "hello",
			runner: func(cmd string) (string, error) {
				return "hello", nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: "",
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name:    "runner return error",
			command: "kb",
			runner: func(cmd string) (string, error) {
				return string(ReadFixtureFile("./testdata/fixtures/host_disk_io", t)), errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostBox(tc.runner, tc.command)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_HostGauge(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected float64
		command  string
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: 60,
			runner: func(cmd string) (string, error) {
				return "60", nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: 0,
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name:     "runner return wrong result",
			expected: 0,
			runner: func(cmd string) (string, error) {
				return "ldfsjsdf", nil
			},
			wantErr: false,
		},
		{
			name:    "runner return error",
			command: "kb",
			runner: func(cmd string) (string, error) {
				return "60", errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostGauge(tc.runner, tc.command)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
func Test_HostBar(t *testing.T) {
	testCases := []struct {
		name     string
		runner   runnerFunc
		expected []int
		command  string
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: []int{40, 50, 60, 70},
			runner: func(cmd string) (string, error) {
				return "40 50 60 70", nil
			},
			wantErr: false,
		},
		{
			name:     "on multiple lines",
			expected: []int{40, 50, 60, 70},
			runner: func(cmd string) (string, error) {
				return `40 50 
						60 70`, nil
			},
			wantErr: false,
		},
		{
			name:     "empty result",
			expected: []int{0},
			runner: func(cmd string) (string, error) {
				return "", nil
			},
			wantErr: false,
		},
		{
			name:     "runner return wrong result",
			expected: []int{0},
			runner: func(cmd string) (string, error) {
				return "ldfsjsdf", nil
			},
			wantErr: false,
		},
		{
			name:    "runner return error",
			command: "kb",
			runner: func(cmd string) (string, error) {
				return "60", errors.New("Error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := HostBar(tc.runner, tc.command)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
