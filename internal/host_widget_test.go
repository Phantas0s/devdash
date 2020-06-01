package internal

import (
	"testing"
	"time"
)

func Test_formatSeconds(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		duration time.Duration
	}{
		{
			name:     "happy case",
			expected: "4h 46m 40s",
			duration: time.Duration(17200210000000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatSeconds(tc.duration)

			if actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
