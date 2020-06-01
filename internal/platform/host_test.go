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
		runner   func(cmd string) (string, error)
		wantErr  bool
	}{
		{
			name:     "happy case",
			expected: 17200210000000,
			runner:   func(cmd string) (string, error) { return "17200.21 59425.48", nil },
			wantErr:  false,
		},
		{
			name:     "Empty result",
			expected: 0,
			runner:   func(cmd string) (string, error) { return "", nil },
			wantErr:  true,
		},
		{
			name:     "Runner return error",
			expected: 0,
			runner:   func(cmd string) (string, error) { return "", errors.New("Error!") },
			wantErr:  true,
		},
		{
			name:     "Runner return impossible number",
			expected: 0,
			runner:   func(cmd string) (string, error) { return "hello", nil },
			wantErr:  true,
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
