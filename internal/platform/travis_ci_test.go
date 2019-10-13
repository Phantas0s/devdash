package platform

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/shuheiktgw/go-travis"
)

func Test_formatBuilds(t *testing.T) {
	testCases := []struct {
		name        string
		fixtureFile string
		expected    [][]string
		limit       int64
	}{
		{
			name: "happy case",
			expected: [][]string{
				{
					"Repository",
					"State",
					"Duration",
					"Finished At",
				},
				{
					"devdash",
					"passed",
					"256",
					"2019-10-11T17:36:48Z",
				},
				{
					"devdash",
					"passed",
					"253",
					"2019-10-10T19:02:03Z",
				},
			},
			fixtureFile: "./testdata/fixtures/travis_table_builds.json",
			limit:       2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tb := []*travis.Build{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, &tb)
			if err != nil {
				t.Error(err)
			}

			actual := formatBuilds(tb, tc.limit)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
