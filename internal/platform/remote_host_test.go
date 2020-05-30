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
