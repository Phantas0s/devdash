package plateform

import (
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
			name: "one filter / contains",
			expected: []*sc.ApiDimensionFilter{
				{
					Dimension:  "query",
					Expression: "test",
					Operator:   "contains",
				},
			},
			dimension: "query",
			filters:   "test",
		},
		{
			name: "one filter / not contains",
			expected: []*sc.ApiDimensionFilter{
				{
					Dimension:  "query",
					Expression: "test",
					Operator:   "notContains",
				},
			},
			dimension: "query",
			filters:   "-test",
		},
		{
			name: "multiple filters",
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
			},
			dimension: "page",
			filters:   "-exclude, include,mobile -useless, query hello",
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
