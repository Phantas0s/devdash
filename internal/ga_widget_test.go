package internal

import (
	"reflect"
	"testing"
	"time"
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

func Test_ExtractDimensions(t *testing.T) {
	testCases := []struct {
		name     string
		expected []string
		options  map[string]string
	}{
		{
			name:     "no dimension in option",
			expected: []string{},
			options:  map[string]string{},
		},
		{
			name:     "one dimension in option",
			expected: []string{"dimension"},
			options:  map[string]string{optionDimensions: "dimension"},
		},
		{
			name:     "multiple dimensions",
			expected: []string{"dimension", "super", "extra"},
			options:  map[string]string{optionDimensions: "dimension,super,extra"},
		},
		{
			name:     "wrong format",
			expected: []string{"dimension super extra"},
			options:  map[string]string{optionDimensions: "dimension super extra"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractDimensions(tc.options)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_ExtractMetric(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		options  map[string]string
	}{
		{
			name:     "no metric in option",
			expected: "sessions",
			options:  map[string]string{},
		},
		{
			name:     "one metric in option",
			expected: "metric",
			options:  map[string]string{optionMetric: "metric"},
		},
		{
			name:     "wrong format with comma delimiters",
			expected: "metric,lala",
			options:  map[string]string{optionMetric: "metric,lala"},
		},
		{
			name:     "wrong format with space delimiters",
			expected: "metric lala",
			options:  map[string]string{optionMetric: "metric lala"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractMetric(tc.options)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_ExtractTimeRange(t *testing.T) {
	testCases := []struct {
		name              string
		expectedStartDate time.Time
		expectedEndDate   time.Time
		options           map[string]string
		base              time.Time
		wantErr           bool
	}{
		{
			name:              "no time, 7 days ago by default",
			base:              time.Date(2020, 5, 5, 0, 0, 0, 0, time.Local),
			expectedStartDate: time.Date(2020, 4, 28, 0, 0, 0, 0, time.Local),
			expectedEndDate:   time.Date(2020, 5, 5, 0, 0, 0, 0, time.Local),
			options:           map[string]string{},
			wantErr:           false,
		},
		{
			name:              "this month",
			base:              time.Date(2020, 5, 5, 0, 0, 0, 0, time.Local),
			expectedStartDate: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local),
			expectedEndDate:   time.Date(2020, 5, 31, 0, 0, 0, 0, time.Local),
			options: map[string]string{
				"start_date": "this_month",
				"end_date":   "this_month",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sd, ed, err := ExtractTimeRange(tc.base, tc.options)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !sd.Equal(tc.expectedStartDate) {
				t.Errorf("Expected start date %v, actual %v", tc.expectedStartDate, sd)
			}
			if tc.wantErr == false && !ed.Equal(tc.expectedEndDate) {
				t.Errorf("Expected end date %v, actual %v", tc.expectedEndDate, ed)
			}
		})
	}
}
