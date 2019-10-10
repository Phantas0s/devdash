package platform

import (
	"encoding/json"
	"reflect"
	"testing"

	sc "google.golang.org/api/webmasters/v3"
)

func Test_filtersFromString(t *testing.T) {
	testCases := []struct {
		name      string
		filters   string
		dimension string
		expected  []*sc.ApiDimensionFilter
	}{
		{
			name:      "one filter / contains",
			dimension: "query",
			filters:   "test",
			expected: []*sc.ApiDimensionFilter{
				{
					Dimension:  "query",
					Expression: "test",
					Operator:   "contains",
				},
			},
		},
		{
			name:      "one filter / not contains",
			dimension: "query",
			filters:   "-test",
			expected: []*sc.ApiDimensionFilter{
				{
					Dimension:  "query",
					Expression: "test",
					Operator:   "notContains",
				},
			},
		},
		{
			name:      "multiple filters",
			dimension: "page",
			filters:   "-exclude, include,*mobile -useless, *query hello, hello halli hallo, *query -hela helo",
			expected: []*sc.ApiDimensionFilter{
				{
					Dimension:  "page",
					Expression: "exclude",
					Operator:   "notContains",
				},
				{
					Dimension:  "page",
					Expression: "include",
					Operator:   "contains",
				},
				{
					Dimension:  "mobile",
					Expression: "useless",
					Operator:   "notContains",
				},
				{
					Dimension:  "query",
					Expression: "hello",
					Operator:   "contains",
				},
				{
					Dimension:  "page",
					Expression: "hello halli hallo",
					Operator:   "contains",
				},
				{
					Dimension:  "query",
					Expression: "hela helo",
					Operator:   "notContains",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := filtersFromString(tc.filters, tc.dimension)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_formatSearchTable(t *testing.T) {
	testCases := []struct {
		name        string
		expected    []SearchConsoleResponse
		fixtureFile string
	}{
		{
			name:        "happy case",
			fixtureFile: "./testdata/fixtures/gsc_query.json",
			expected: []SearchConsoleResponse{
				{
					Dimension:   "php compare datetime",
					Clicks:      71,
					Impressions: 182,
					Ctr:         0.3901098901098901,
					Position:    1.7142857142857144,
				},
				{
					Dimension:   "php datetime compare",
					Clicks:      41,
					Impressions: 110,
					Ctr:         0.37272727272727274,
					Position:    1.518181818181818,
				},
				{
					Dimension:   "php compare dates",
					Clicks:      34,
					Impressions: 439,
					Ctr:         0.0774487471526196,
					Position:    3.5261958997722096,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &sc.SearchAnalyticsQueryResponse{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, r)
			if err != nil {
				t.Error(err)
			}
			actual := formatSearchTable(r.Rows)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
