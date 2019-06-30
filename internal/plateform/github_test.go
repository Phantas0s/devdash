package plateform

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/github"
)

// StartDate is normally the current date
func Test_CommitCount(t *testing.T) {
	testCases := []struct {
		name        string
		fixtureFile string
		expectedDim []string
		expectedVal []int
		startWeek   int64
		endWeek     int64
		startDate   time.Time
	}{
		{
			name:        "format commit counts 6_weeks_ago to today",
			expectedVal: []int{3, 31, 6, 0, 0, 0},
			expectedDim: []string{"05-26", "06-02", "06-09", "06-16", "06-23", "06-30"},
			fixtureFile: "./testdata/fixtures/github_commit_count.json",
			startWeek:   6,                                                 // 6_weeks_ago
			endWeek:     0,                                                 // today
			startDate:   time.Date(2019, 06, 30, 00, 00, 00, 00, time.UTC), // Sunday 30 June 2019
		},
		{
			name:        "format commit counts 6_weeks_ago to today, beginning by Thursday",
			expectedVal: []int{3, 31, 6, 0, 0, 0},
			expectedDim: []string{"05-19", "05-26", "06-02", "06-09", "06-16", "06-23"},
			fixtureFile: "./testdata/fixtures/github_commit_count.json",
			startWeek:   6,                                                 // 6_weeks_ago
			endWeek:     0,                                                 // today
			startDate:   time.Date(2019, 06, 27, 00, 00, 00, 00, time.UTC), // Thursday 27 June 2019
		},
		{
			name:        "format commit counts 5_weeks_ago to 1_weeks_ago",
			expectedVal: []int{31, 6, 0, 0},
			expectedDim: []string{"06-02", "06-09", "06-16", "06-23"},
			fixtureFile: "./testdata/fixtures/github_commit_count.json",
			startWeek:   5,                                                 // 5_weeks_ago
			endWeek:     1,                                                 // 1_weeks_ago
			startDate:   time.Date(2019, 06, 30, 00, 00, 00, 00, time.UTC), // Sunday 30 June 2019
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			part := &github.RepositoryParticipation{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, part)
			if err != nil {
				t.Error(err)
			}

			dim, val := formatCommitCounts(part, tc.startWeek, tc.endWeek, tc.startDate)

			if !reflect.DeepEqual(dim, tc.expectedDim) {
				t.Errorf("Expected %v, actual %v", tc.expectedDim, dim)
			}

			if !reflect.DeepEqual(val, tc.expectedVal) {
				t.Errorf("Expected %v, actual %v", tc.expectedVal, val)
			}
		})
	}
}
