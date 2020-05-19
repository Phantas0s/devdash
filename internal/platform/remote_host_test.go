package platform

import (
	"reflect"
	"testing"
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
				{
					"data1", "data2", "data3", "data4",
				},
				{
					"data5", "data6", "data7", "data8",
				},
			},
			headers: []string{"dataHeader1", "dataHeader2", "dataHeader3", "dataHeader4"},
			data:    "data1,data2,data3,data4,data5,data6,data7,data8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatToTable(tc.headers, tc.data)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_ConvertUnit(t *testing.T) {
	testCases := []struct {
		name     string
		expected []int
		input    []int
		base     string
		to       string
	}{
		{
			name:     "kb->mb",
			expected: []int{1, 2, 3},
			input:    []int{1024, 2048, 4000},
			base:     "kb",
			to:       "mb",
		},
		{
			name:     "kb->gb",
			expected: []int{10, 6},
			input:    []int{10485760, 6291456},
			base:     "kb",
			to:       "gb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ConvertBinUnit(tc.input, tc.base, tc.to)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
