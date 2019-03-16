package plateform

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	ga "google.golang.org/api/analyticsreporting/v4"
)

func Test_Format(t *testing.T) {
	testCases := []struct {
		name        string
		expectedVal []int
		fixtureFile string
		expectedDim []string
		formater    func([]string) string
		wantErr     bool
	}{
		{
			name:        "format users",
			expectedVal: []int{370, 414, 387, 202},
			expectedDim: []string{"02-05", "02-06", "02-07", "02-08"},
			fixtureFile: "./testdata/fixtures/ga_users.json",
			formater:    func(dim []string) string { return dim[0] + "-" + dim[1] },
			wantErr:     false,
		},
		{
			name:        "custom bar metric with aggregation",
			expectedVal: []int{35, 40, 29},
			expectedDim: []string{"03-01", "03-02", "03-03"},
			fixtureFile: "./testdata/fixtures/ga_bar_metric_aggregations.json",
			formater:    func(dim []string) string { return dim[0] + "-" + dim[1] },
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret := &ga.GetReportsResponse{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, ret)
			if err != nil {
				t.Error(err)
			}

			dim, val, err := formatBar(ret.Reports, tc.formater)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(dim, tc.expectedDim) {
				t.Errorf("Expected %v, actual %v", tc.expectedDim, dim)
			}

			if tc.wantErr == false && !reflect.DeepEqual(val, tc.expectedVal) {
				t.Errorf("Expected %v, actual %v", tc.expectedVal, val)
			}
		})
	}
}

// TODO simplify the test (less data)... but keep this two month period
func Test_FormatNewReturning(t *testing.T) {
	testCases := []struct {
		name        string
		expectedVal []int
		fixtureFile string
		expectedDim []string
		formater    func([]string) string
		wantErr     bool
	}{
		{
			name: "format new vs returning",
			expectedVal: []int{
				141,
				126,
				302,
				364,
				326,
				329,
				269,
				176,
				120,
				291,
				316,
				364,
				326,
				273,
				164,
				105,
				58,
				40,
				72,
				95,
				96,
				93,
				66,
				35,
				45,
				86,
				78,
				83,
				86,
				80,
				44,
				39,
			},
			expectedDim: []string{
				"01-26",
				"01-27",
				"01-28",
				"01-29",
				"01-30",
				"01-31",
				"02-01",
				"02-02",
				"02-03",
				"02-04",
				"02-05",
				"02-06",
				"02-07",
				"02-08",
				"02-09",
				"02-10",
			},
			fixtureFile: "./testdata/fixtures/ga_new_returning.json",
			formater:    func(dim []string) string { return dim[1] + "-" + dim[2] },
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret := &ga.GetReportsResponse{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, ret)
			if err != nil {
				t.Error(err)
			}

			dim, val, err := formatNewReturning(ret.Reports, tc.formater)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(dim, tc.expectedDim) {
				t.Errorf("Expected %v, actual %v", tc.expectedDim, dim)
			}

			if tc.wantErr == false && !reflect.DeepEqual(val, tc.expectedVal) {
				t.Errorf("Expected %v, actual %v", tc.expectedVal, val)
			}
		})
	}
}

func Test_FormatBarReturning(t *testing.T) {
	testCases := []struct {
		name        string
		expectedVal []int
		fixtureFile string
		expectedDim []string
		formater    func([]string) string
		wantErr     bool
	}{
		{
			name:        "format new vs returning",
			expectedVal: []int{87, 78, 73, 58, 54, 20},
			expectedDim: []string{"03-04", "03-05", "03-06", "03-07", "03-08", "03-09"},
			fixtureFile: "./testdata/fixtures/ga_users_returning.json",
			formater:    func(dim []string) string { return dim[0] + "-" + dim[1] },
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret := &ga.GetReportsResponse{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, ret)
			if err != nil {
				t.Error(err)
			}

			dim, val, err := formatBarReturning(ret.Reports, tc.formater)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(dim, tc.expectedDim) {
				t.Errorf("Expected %v, actual %v", tc.expectedDim, dim)
			}

			if tc.wantErr == false && !reflect.DeepEqual(val, tc.expectedVal) {
				t.Errorf("Expected %v, actual %v", tc.expectedVal, val)
			}
		})
	}
}

func ReadFixtureFile(file string, t *testing.T) (data []byte) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Error(errors.Wrapf(err, "can't read file %s", file))
	}

	return data
}

func Test_mapMetrics(t *testing.T) {
	testCases := []struct {
		name     string
		m        []string
		expected []*ga.Metric
	}{
		{
			name: "happy case",
			m: []string{
				"sessions",
				"page_views",
				"entrances",
				"unique_page_views",
			},
			expected: []*ga.Metric{
				{Expression: "ga:sessions"},
				{Expression: "ga:pageViews"},
				{Expression: "ga:entrances"},
				{Expression: "ga:uniquePageviews"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapMetrics(tc.m)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_mapDimension(t *testing.T) {
	testCases := []struct {
		name     string
		d        string
		expected string
	}{
		{
			name:     "alias",
			d:        "page_path",
			expected: "ga:pagePath",
		},
		{
			name:     "ga name",
			d:        "ga:pagePath",
			expected: "ga:pagePath",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapDimension(tc.d)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_mapHeaders(t *testing.T) {
	testCases := []struct {
		name     string
		m        []string
		expected []string
		el       string
	}{
		{
			name: "happy case",
			el:   "Pages",
			expected: []string{
				"Pages",
				"Sessions",
				"Page Views",
				"Entrances",
				"Unique Page Views",
				"someRandomExpr",
				"hey",
			},
			m: []string{
				"sessions",
				"page_views",
				"entrances",
				"unique_page_views",
				"ga:someRandomExpr",
				" hey",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapHeaders(tc.el, tc.m)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_mapOrderBy(t *testing.T) {
	testCases := []struct {
		name     string
		m        []string
		expected []*ga.OrderBy
	}{
		{
			name: "happy case",
			m: []string{
				"sessions asc",
				"page_views desc",
				"unique_page_views",
				"ga:uniquePageviews",
			},
			expected: []*ga.OrderBy{
				&ga.OrderBy{
					FieldName: "ga:sessions",
					SortOrder: "ASCENDING",
				},
				&ga.OrderBy{
					FieldName: "ga:pageViews",
					SortOrder: "DESCENDING",
				},
				&ga.OrderBy{
					FieldName: "ga:uniquePageviews",
					SortOrder: "DESCENDING",
				},
				&ga.OrderBy{
					FieldName: "ga:uniquePageviews",
					SortOrder: "DESCENDING",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapOrderBy(tc.m)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
