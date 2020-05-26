package gokit

import (
	"reflect"
	"testing"
)

func Test_ConvertUnit(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		input    float64
		base     string
		to       string
	}{
		{
			name:     "kb->mb",
			expected: 1,
			input:    1024,
			base:     "kb",
			to:       "mb",
		},
		{
			name:     "kb->gb",
			expected: 10,
			input:    10485760,
			base:     "kb",
			to:       "gb",
		},
		{
			name:     "kb->kb",
			expected: 10,
			input:    10,
			base:     "kb",
			to:       "kb",
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
