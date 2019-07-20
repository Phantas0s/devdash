package plateform

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v27/github"
)

// StartDate is normally the current date
func Test_formatCountStars(t *testing.T) {
	testCases := []struct {
		name        string
		fixtureFile string
		expectedDim []string
		expectedVal []int
	}{
		// TODO wrong result...
		{
			name:        "happy case",
			expectedVal: []int{36, 21, 23, 3},
			expectedDim: []string{"05-28", "05-29", "05-30", "06-10"},
			fixtureFile: "./testdata/fixtures/github_start_count.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sg1 := []*github.Stargazer{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, &sg1)
			if err != nil {
				t.Error(err)
			}

			dim, val := formatCountStars(sg1)

			if !reflect.DeepEqual(dim, tc.expectedDim) {
				t.Errorf("Expected %v, actual %v", tc.expectedDim, dim)
			}

			if !reflect.DeepEqual(val, tc.expectedVal) {
				t.Errorf("Expected %v, actual %v", tc.expectedVal, val)
			}
		})
	}
}

// StartDate is normally the current date
func Test_formatCommitCount(t *testing.T) {
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

func Test_FomatListPullRequest(t *testing.T) {
	testCases := []struct {
		name        string
		expected    [][]string
		fixtureFile string
		limit       int
	}{
		{
			name: "happy case",
			expected: [][]string{
				{
					"title",
					"state",
					"created at",
					"merged",
					"commits",
				},
				{
					"super pull request",
					"closed",
					"2018-10-19 21:12:25 +0000 UTC",
					"unknown",
					"unknown",
				},
			},
			fixtureFile: "./testdata/fixtures/github_list_pull_request.json",
			limit:       1000000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gpr := []*github.PullRequest{}
			fixtures := ReadFixtureFile(tc.fixtureFile, t)
			err := json.Unmarshal(fixtures, &gpr)
			if err != nil {
				t.Error(err)
			}

			actual := formatListPullRequests(gpr, tc.limit)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
