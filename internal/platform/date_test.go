package platform

import (
	"testing"
	"time"
)

func Test_ConvertStartDate(t *testing.T) {
	testCases := []struct {
		name              string
		startDate         string
		endDate           string
		expectedStartDate string
		expectedEndDate   string
		baseYear          int
		baseMonth         time.Month
		baseDay           int
		wantErr           bool
	}{
		{
			name:              "today",
			startDate:         "today",
			endDate:           "today",
			expectedStartDate: "2019-03-27",
			expectedEndDate:   "2019-03-27",
			baseYear:          2019,
			baseMonth:         03,
			baseDay:           27,
			wantErr:           false,
		},
		{
			name:              "yesterday",
			startDate:         "yesterday",
			endDate:           "yesterday",
			expectedStartDate: "2019-03-26",
			expectedEndDate:   "2019-03-26",
			baseYear:          2019,
			baseMonth:         03,
			baseDay:           27,
			wantErr:           false,
		},
		{
			name:              "previous days",
			startDate:         "7_days_ago",
			endDate:           "7_days_ago",
			expectedStartDate: "2019-03-20",
			expectedEndDate:   "2019-03-20",
			baseYear:          2019,
			baseMonth:         03,
			baseDay:           27,
			wantErr:           false,
		},
		{
			name:              "this week",
			startDate:         "this_week",
			endDate:           "this_week",
			expectedStartDate: "2019-03-11",
			expectedEndDate:   "2019-03-17",
			baseYear:          2019,
			baseMonth:         03,
			baseDay:           14,
			wantErr:           false,
		},
		{
			name:              "previous weeks",
			startDate:         "2_weeks_ago",
			endDate:           "2_weeks_ago",
			expectedStartDate: "2019-03-11",
			expectedEndDate:   "2019-03-17",
			baseYear:          2019,
			baseMonth:         03,
			baseDay:           27,
			wantErr:           false,
		},
		{
			name:              "this month",
			startDate:         "this_month",
			endDate:           "this_month",
			expectedStartDate: "2019-01-01",
			expectedEndDate:   "2019-01-31",
			baseYear:          2019,
			baseMonth:         01,
			baseDay:           27,
			wantErr:           false,
		},
		{
			name:              "previous months",
			startDate:         "2_months_ago",
			endDate:           "2_months_ago",
			expectedStartDate: "2019-01-01",
			expectedEndDate:   "2019-01-31",
			baseYear:          2019,
			baseMonth:         03,
			baseDay:           27,
			wantErr:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			base := time.Date(tc.baseYear, tc.baseMonth, tc.baseDay, 00, 00, 00, 00, time.UTC)
			start, end, err := ConvertDates(base, tc.startDate, tc.endDate)
			if (err != nil) != tc.wantErr {
				t.Errorf("Error '%v' even if wantErr is %t", err, tc.wantErr)
				return
			}

			if tc.wantErr == false && tc.expectedStartDate != start.Format("2006-01-02") {
				t.Errorf("Expected start date%v, actual %v", tc.expectedStartDate, start)
			}

			if tc.wantErr == false && tc.expectedEndDate != end.Format("2006-01-02") {
				t.Errorf("Expected end date%v, actual %v", tc.expectedEndDate, end)
			}
		})
	}
}
