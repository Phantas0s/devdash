package internal

import (
	"reflect"
	"testing"
)

func Test_formatTable(t *testing.T) {
	testCases := []struct {
		name      string
		expected  [][]string
		rowLimit  int64
		dim       []string
		val       [][]string
		charLimit int64
		headers   []string
	}{
		{
			name:      "happy case",
			rowLimit:  2,
			charLimit: 3,
			dim:       []string{"/php/ ", "     /go/", " /c++/"},
			val: [][]string{
				{"12"},
				{"123"},
				{"2"},
			},
			headers: []string{"Page", "Sessions"},
			expected: [][]string{
				{"Page", "Sessions"},
				{"/ph", "12"},
				{"/go", "123"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatTable(tc.rowLimit, tc.dim, tc.val, tc.charLimit, tc.headers)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
