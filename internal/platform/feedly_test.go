package platform

import (
	"reflect"
	"testing"
)

func Test_extractSubscribers(t *testing.T) {
	testCases := []struct {
		name        string
		expected    string
		fixtureFile string
		wantErr     bool
	}{
		{
			name:        "happy case",
			expected:    "100",
			fixtureFile: "./testdata/fixtures/feedly_search.json",
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixtures := ReadFixtureFile(tc.fixtureFile, t)

			actual, err := extractSubscribers(fixtures)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
