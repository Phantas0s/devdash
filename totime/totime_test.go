package totime

import (
	"testing"
	"time"
)

const yyyymmdd = "2006-01-02"

func Test_ThisWeek(t *testing.T) {
	testCases := []struct {
		name      string
		startDate string
		endDate   string
		current   time.Time
		wantErr   bool
	}{
		{
			name:      "monday",
			startDate: "2019-03-11",
			endDate:   "2019-03-17",
			current:   time.Date(2019, 03, 11, 00, 00, 00, 00, time.UTC),
			wantErr:   false,
		},
		{
			name:      "another day",
			startDate: "2019-03-11",
			endDate:   "2019-03-17",
			current:   time.Date(2019, 03, 14, 00, 00, 00, 00, time.UTC),
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := ThisWeek(tc.current)

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
		current      time.Time
		weekInFuture int
		wantErr      bool
	}{
		{
			name:         "monday",
			startDate:    "2019-03-18",
			endDate:      "2019-03-24",
			current:      time.Date(2019, 03, 11, 00, 00, 00, 00, time.UTC),
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "middle of week",
			startDate:    "2019-03-18",
			endDate:      "2019-03-24",
			current:      time.Date(2019, 03, 13, 00, 00, 00, 00, time.UTC),
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "end of month",
			startDate:    "2019-04-01",
			endDate:      "2019-04-07",
			current:      time.Date(2019, 03, 28, 00, 00, 00, 00, time.UTC),
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "end of year",
			startDate:    "2018-12-31",
			endDate:      "2019-01-06",
			current:      time.Date(2018, 12, 25, 00, 00, 00, 00, time.UTC),
			weekInFuture: 1,
			wantErr:      false,
		},
		{
			name:         "end of year multiple weeks",
			startDate:    "2019-01-14",
			endDate:      "2019-01-20",
			current:      time.Date(2018, 12, 25, 00, 00, 00, 00, time.UTC),
			weekInFuture: 3,
			wantErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := NextWeeks(tc.current, tc.weekInFuture)

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
		name       string
		startDate  string
		endDate    string
		current    time.Time
		weekInPast int
		wantErr    bool
	}{
		{
			name:       "monday",
			startDate:  "2019-03-04",
			endDate:    "2019-03-10",
			current:    time.Date(2019, 03, 11, 00, 00, 00, 00, time.UTC),
			weekInPast: 1,
			wantErr:    false,
		},
		{
			name:       "middle of week",
			startDate:  "2019-03-04",
			endDate:    "2019-03-10",
			current:    time.Date(2019, 03, 13, 00, 00, 00, 00, time.UTC),
			weekInPast: 1,
			wantErr:    false,
		},
		{
			name:       "beginning of month",
			startDate:  "2019-03-25",
			endDate:    "2019-03-31",
			current:    time.Date(2019, 04, 03, 00, 00, 00, 00, time.UTC),
			weekInPast: 1,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := PrevWeeks(tc.current, tc.weekInPast)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}

func Test_ThisMonth(t *testing.T) {
	testCases := []struct {
		name      string
		startDate string
		endDate   string
		current   time.Time
		wantErr   bool
	}{
		{
			name:      "monday",
			startDate: "2019-03-01",
			endDate:   "2019-03-31",
			current:   time.Date(2019, 03, 11, 00, 00, 00, 00, time.UTC),
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := ThisMonth(tc.current)

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
		current    time.Time
		countMonth int
		wantErr    bool
	}{
		{
			name:       "february",
			startDate:  "2017-02-01",
			endDate:    "2017-02-28",
			current:    time.Date(2017, 03, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "february bisextile",
			startDate:  "2008-02-01",
			endDate:    "2008-02-29",
			current:    time.Date(2008, 03, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "beginning of the year",
			startDate:  "2017-12-01",
			endDate:    "2017-12-31",
			current:    time.Date(2018, 01, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "end of the year",
			startDate:  "2018-11-01",
			endDate:    "2018-11-30",
			current:    time.Date(2018, 12, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := PrevMonths(tc.current, tc.countMonth)

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
		current    time.Time
		countMonth int
		wantErr    bool
	}{
		{
			name:       "february",
			startDate:  "2017-02-01",
			endDate:    "2017-02-28",
			current:    time.Date(2017, 01, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "february bisextile",
			startDate:  "2008-02-01",
			endDate:    "2008-02-29",
			current:    time.Date(2008, 01, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "end of the year",
			startDate:  "2018-01-01",
			endDate:    "2018-01-31",
			current:    time.Date(2017, 12, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 1,
			wantErr:    false,
		},
		{
			name:       "this month",
			startDate:  "2018-01-01",
			endDate:    "2018-01-31",
			current:    time.Date(2018, 01, 18, 00, 00, 00, 00, time.UTC),
			countMonth: 0,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := NextMonths(tc.current, tc.countMonth)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}

func Test_ThisYear(t *testing.T) {
	testCases := []struct {
		name      string
		startDate string
		endDate   string
		current   time.Time
		wantErr   bool
	}{
		{
			name:      "monday",
			startDate: "2019-01-01",
			endDate:   "2019-12-31",
			current:   time.Date(2019, 03, 11, 00, 00, 00, 00, time.UTC),
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := ThisYear(tc.current)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}
func Test_NextYears(t *testing.T) {
	testCases := []struct {
		name       string
		startDate  string
		endDate    string
		current    time.Time
		countYears int
		wantErr    bool
	}{
		{
			name:       "1 year in future",
			startDate:  "2019-01-01",
			endDate:    "2019-12-31",
			current:    time.Date(2018, 01, 18, 00, 00, 00, 00, time.UTC),
			countYears: 1,
			wantErr:    false,
		},
		{
			name:       "5 year in future",
			startDate:  "2023-01-01",
			endDate:    "2023-12-31",
			current:    time.Date(2018, 01, 18, 00, 00, 00, 00, time.UTC),
			countYears: 5,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := NextYears(tc.current, tc.countYears)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}

func Test_PrevYears(t *testing.T) {
	testCases := []struct {
		name       string
		startDate  string
		endDate    string
		current    time.Time
		countYears int
		wantErr    bool
	}{
		{
			name:       "1 year in past",
			startDate:  "2017-01-01",
			endDate:    "2017-12-31",
			countYears: 1,
			current:    time.Date(2018, 01, 18, 00, 00, 00, 00, time.UTC),
		},
		{
			name:       "5 years in past",
			startDate:  "2013-01-01",
			endDate:    "2013-12-31",
			countYears: 5,
			current:    time.Date(2018, 01, 18, 00, 00, 00, 00, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, endDate := PrevYears(tc.current, tc.countYears)

			if tc.wantErr == false && tc.startDate != startDate.Format(yyyymmdd) {
				t.Errorf("Expected startDate %v, actual %v", tc.startDate, startDate)
			}

			if tc.wantErr == false && tc.endDate != endDate.Format(yyyymmdd) {
				t.Errorf("Expected endDate %v, actual %v", tc.endDate, endDate)
			}
		})
	}
}
