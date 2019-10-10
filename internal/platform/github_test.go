package platform

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v28/github"
)

// StartDate is normally the current date
func Test_formatCountStars(t *testing.T) {
	testCases := []struct {
		name        string
		fixtureFile string
		expectedDim []string
		expectedVal []int
		startDate   time.Time
		endDate     time.Time
		timeLayout  string
	}{
		{
			name:        "happy case",
			expectedVal: []int{36, 21, 23, 3},
			expectedDim: []string{"05-28", "05-29", "05-30", "06-10"},
			startDate:   time.Date(2019, 05, 01, 00, 00, 00, 00, time.UTC),
			endDate:     time.Date(2019, 07, 01, 00, 00, 00, 00, time.UTC),
			fixtureFile: "./testdata/fixtures/github_start_count.json",
			timeLayout:  "01-02",
		},
		{
			name:        "time restriction",
			expectedVal: []int{21, 23},
			expectedDim: []string{"05-29", "05-30"},
			startDate:   time.Date(2019, 05, 29, 00, 00, 00, 00, time.UTC),
			endDate:     time.Date(2019, 06, 01, 00, 00, 00, 00, time.UTC),
			fixtureFile: "./testdata/fixtures/github_start_count.json",
			timeLayout:  "01-02",
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

			dim, val := formatCountStars(sg1, tc.timeLayout, tc.startDate, tc.endDate, false)

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
		scope       string
	}{
		{
			name:        "format commit counts 6_weeks_ago to today",
			expectedVal: []int{3, 31, 6, 0, 0, 0},
			expectedDim: []string{"05-26", "06-02", "06-09", "06-16", "06-23", "06-30"},
			fixtureFile: "./testdata/fixtures/github_commit_count.json",
			startWeek:   6,                                                 // 6_weeks_ago
			endWeek:     0,                                                 // today
			startDate:   time.Date(2019, 06, 30, 00, 00, 00, 00, time.UTC), // Sunday 30 June 2019
			scope:       githubScopeOwner,
		},
		{
			name:        "format commit counts 6_weeks_ago to today, beginning by Thursday",
			expectedVal: []int{3, 31, 6, 0, 0, 0},
			expectedDim: []string{"05-19", "05-26", "06-02", "06-09", "06-16", "06-23"},
			fixtureFile: "./testdata/fixtures/github_commit_count.json",
			startWeek:   6,                                                 // 6_weeks_ago
			endWeek:     0,                                                 // today
			startDate:   time.Date(2019, 06, 27, 00, 00, 00, 00, time.UTC), // Thursday 27 June 2019
			scope:       githubScopeOwner,
		},
		{
			name:        "format commit counts 5_weeks_ago to 1_weeks_ago",
			expectedVal: []int{31, 6, 0, 0},
			expectedDim: []string{"06-02", "06-09", "06-16", "06-23"},
			fixtureFile: "./testdata/fixtures/github_commit_count.json",
			startWeek:   5,                                                 // 5_weeks_ago
			endWeek:     1,                                                 // 1_weeks_ago
			startDate:   time.Date(2019, 06, 30, 00, 00, 00, 00, time.UTC), // Sunday 30 June 2019
			scope:       githubScopeOwner,
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

			p := part.Owner
			if tc.scope == githubScopeAll {
				p = part.All
			}

			dim, val := formatCountCommits(p, tc.startWeek, tc.endWeek, tc.startDate)

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

func Test_fillMissingDays(t *testing.T) {
	testCases := []struct {
		name           string
		expected       []time.Time
		dim            []time.Time
		values         []int
		expectedValues []int
	}{
		{
			name: "happy case",
			expected: []time.Time{
				time.Date(2019, time.January, 01, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 02, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 03, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 04, 0, 0, 0, 0, time.Now().Location()),
			},
			dim: []time.Time{
				time.Date(2019, time.January, 01, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 04, 0, 0, 0, 0, time.Now().Location()),
			},
			values: []int{
				12,
				15,
			},
			expectedValues: []int{
				12,
				0,
				0,
				15,
			},
		},
		{
			name: "more complex case",
			expected: []time.Time{
				time.Date(2019, time.January, 01, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 02, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 03, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 04, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 05, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 06, 0, 0, 0, 0, time.Now().Location()),
			},
			dim: []time.Time{
				time.Date(2019, time.January, 01, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 04, 0, 0, 0, 0, time.Now().Location()),
				time.Date(2019, time.January, 06, 0, 0, 0, 0, time.Now().Location()),
			},
			values: []int{
				12,
				15,
				18,
			},
			expectedValues: []int{
				12,
				0,
				0,
				15,
				0,
				18,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, val := fillMissingDays(tc.dim, tc.values)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected dates %v, actual %v", tc.expected, actual)
			}

			if !reflect.DeepEqual(val, tc.expectedValues) {
				t.Errorf("Expected values %v, actual %v", tc.expectedValues, val)
			}
		})
	}
}
