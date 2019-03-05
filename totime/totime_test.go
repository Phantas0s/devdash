package totime

import (
	"testing"
	"time"
)

const yyyymmdd = "2006-01-02"

func Test_ThisMonth(t *testing.T) {
	testCases := []struct {
		name      string
		startDate string
		endDate   string
		year      int
		month     time.Month
		day       int
		wantErr   bool
	}{
		{
			name:      "monday",
			startDate: "2019-03-01",
			endDate:   "2019-03-31",
			year:      2019,
			month:     03,
			day:       11,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)
			startDate, endDate := ThisMonth(current)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}
func Test_PrevMonths(t *testing.T) {
	testCases := []struct {
		name       string
		startDate  string
		endDate    string
		year       int
		month      time.Month
		day        int
		countMonth int
		wantErr    bool
	}{
		{
			name:       "february",
			startDate:  "2017-02-01",
			endDate:    "2017-02-28",
			year:       2017,
			month:      03,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "february bisextile",
			startDate:  "2008-02-01",
			endDate:    "2008-02-29",
			year:       2008,
			month:      03,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "beginning of the year",
			startDate:  "2017-12-01",
			endDate:    "2017-12-31",
			year:       2018,
			month:      01,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "end of the year",
			startDate:  "2018-11-01",
			endDate:    "2018-11-30",
			year:       2018,
			month:      12,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current := time.Date(tc.year, tc.month, tc.day, 00, 00, 00, 00, time.UTC)
			startDate, endDate := PrevMonths(current, tc.countMonth)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}

func Test_NextMonths(t *testing.T) {
	testCases := []struct {
		name       string
		startDate  string
		endDate    string
		year       int
		month      time.Month
		day        int
		countMonth int
		wantErr    bool
	}{
		{
			name:       "february",
			startDate:  "2017-02-01",
			endDate:    "2017-02-28",
			year:       2017,
			month:      01,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "february bisextile",
			startDate:  "2008-02-01",
			endDate:    "2008-02-29",
			year:       2008,
			month:      01,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "end of the year",
			startDate:  "2018-01-01",
			endDate:    "2018-01-31",
			year:       2017,
			month:      12,
			day:        18,
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "this month",
			startDate:  "2018-01-01",
			endDate:    "2018-01-31",
			year:       2018,
			month:      01,
			day:        18,
			countMonth: 0,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current := time.Date(tc.year, tc.month, tc.day, 00, 00, 00, 00, time.UTC)
			startDate, endDate := NextMonths(current, tc.countMonth)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}

func Test_ThisWeek(t *testing.T) {
	testCases := []struct {
		name      string
		startDate string
		endDate   string
		year      int
		month     time.Month
		day       int
		wantErr   bool
	}{
		{
			name:      "monday",
			startDate: "2019-03-11",
			endDate:   "2019-03-17",
			year:      2019,
			month:     03,
			day:       11,
			wantErr:   false,
		},
		{
			name:      "another day",
			startDate: "2019-03-11",
			endDate:   "2019-03-17",
			year:      2019,
			month:     03,
			day:       14,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current := time.Date(tc.year, tc.month, tc.day, 00, 00, 00, 00, time.UTC)
			startDate, endDate := ThisWeek(current)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}
func Test_NextWeeks(t *testing.T) {
	testCases := []struct {
		name         string
		startDate    string
		endDate      string
		year         int
		month        time.Month
		day          int
		weekInFuture int
		wantErr      bool
	}{
		{
			name:         "monday",
			startDate:    "2019-03-18",
			endDate:      "2019-03-24",
			year:         2019,
			month:        03,
			day:          11,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "middle of week",
			startDate:    "2019-03-18",
			endDate:      "2019-03-24",
			year:         2019,
			month:        03,
			day:          13,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "end of month",
			startDate:    "2019-04-01",
			endDate:      "2019-04-07",
			year:         2019,
			month:        03,
			day:          28,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "end of year",
			startDate:    "2018-12-31",
			endDate:      "2019-01-06",
			year:         2018,
			month:        12,
			day:          25,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "end of year multiple weeks",
			startDate:    "2019-01-14",
			endDate:      "2019-01-20",
			year:         2018,
			month:        12,
			day:          25,
			weekInFuture: 3,
			wantErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current := time.Date(tc.year, tc.month, tc.day, 00, 00, 00, 00, time.UTC)
			startDate, endDate := NextWeeks(current, tc.weekInFuture)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}

func Test_PrevWeeks(t *testing.T) {
	testCases := []struct {
		name         string
		startDate    string
		endDate      string
		year         int
		month        time.Month
		day          int
		weekInFuture int
		wantErr      bool
	}{
		{
			name:         "monday",
			startDate:    "2019-03-04",
			endDate:      "2019-03-10",
			year:         2019,
			month:        03,
			day:          11,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "middle of week",
			startDate:    "2019-03-04",
			endDate:      "2019-03-10",
			year:         2019,
			month:        03,
			day:          13,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "beginning of month",
			startDate:    "2019-03-25",
			endDate:      "2019-03-31",
			year:         2019,
			month:        04,
			day:          03,
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "begining of year",
			startDate:    "2018-12-24",
			endDate:      "2018-12-30",
			year:         2019,
			month:        01,
			day:          02,
			weekInFuture: 1,
			wantErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current := time.Date(tc.year, tc.month, tc.day, 00, 00, 00, 00, time.UTC)
			startDate, endDate := PrevWeeks(current, tc.weekInFuture)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}
