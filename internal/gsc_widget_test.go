package internal

import (
	"reflect"
	"testing"

	"github.com/Phantas0s/devdash/internal/platform"
)

func Test_formatText(t *testing.T) {
	testCases := []struct {
		name       string
		trimPrefix string
		expected   [][]string
		table      [][]string
		charLimit  int
	}{
		{
			name:       "happy case",
			trimPrefix: "https://web-techno.net",
			expected: [][]string{
				{"query", "clicks", "impressions", "ctr", "position"},
				{"This", "13", "230", "1"},
				{"Woww", "13", "230", "1"},
				{"/Wow", "13", "230", "1"},
			},
			charLimit: 4,
			table: [][]string{
				{"query", "clicks", "impressions", "ctr", "position"},
				{"This is the first element", "13", "230", "1"},
				{"Wowwoupi is the second element", "13", "230", "1"},
				{"https://web-techno.net/Wowwoupi/youp", "13", "230", "1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatText(tc.table, tc.charLimit, tc.trimPrefix)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_formatNumerics(t *testing.T) {
	testCases := []struct {
		name      string
		expected  [][]string
		results   []platform.SearchConsoleResponse
		metrics   []string
		dimension string
	}{
		{
			name: "happy case",
			expected: [][]string{
				{"Query", "clicks", "impressions", "ctr", "position"},
				{"this is a query", "10.2", "10.1", "2.37%", "1.23"},
			},
			results: []platform.SearchConsoleResponse{
				{
					Dimension:   "this is a query",
					Clicks:      10.2,
					Impressions: 10.1,
					Ctr:         0.02365,
					Position:    1.234,
				},
			},
			metrics:   []string{"clicks", "impressions", "ctr", "position"},
			dimension: "query",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatNumerics(tc.results, tc.dimension, tc.metrics)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
