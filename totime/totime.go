// totime is a simple package which return the start date and end date of a precise period of time.
// The period of time can be this week, this month, previous week, previous month, the month in two months...
//
// Examples:
// We are in January 2019, ThisMonth will return startDate = 2019-01-01, endDate = 2019-01-31
// We are in January 2019, NextMonth will return startDate = 2019-02-01, endDate = 2019-02-28
// We are in January 2019, NextMonth with count of 2 will return startDate = 2019-03-01, endDate = 2019-03-31

package totime

import (
	"time"
)

func ThisWeek(base time.Time) (startDate time.Time, endDate time.Time) {
	startDate = time.Time{}

	// 1 = Monday
	weekDay := int(base.Weekday())
	startDate = base.AddDate(0, 0, -(weekDay - 1))
	if weekDay == 0 {
		startDate = base.AddDate(0, 0, -(weekDay + 1))
	}

	endDate = startDate.AddDate(0, 0, 6)

	return
}

func PrevWeeks(base time.Time, count int) (startDate time.Time, endDate time.Time) {
	startDate = time.Time{}

	// 1 = Monday
	weekDay := int(base.Weekday())
	startDate = base.AddDate(0, 0, (-(weekDay - 1) - (7 * count)))
	if weekDay == 0 {
		startDate = base.AddDate(0, 0, (-(weekDay + 1) - (7 * count)))
	}

	endDate = startDate.AddDate(0, 0, 6)

	return
}

func NextWeeks(base time.Time, count int) (startDate time.Time, endDate time.Time) {
	// 1 = Monday
	weekDay := int(base.Weekday())
	startDate = base.AddDate(0, 0, (-(weekDay - 1) + (7 * count)))
	if weekDay == 0 {
		startDate = base.AddDate(0, 0, (-(weekDay + 1) + (7 * count)))
	}

	endDate = startDate.AddDate(0, 0, 6)

	return
}

func ThisMonth(base time.Time) (startDate time.Time, endDate time.Time) {
	currentLocation := base.Location()

	startDate = time.Date(base.Year(), base.Month(), 1, 0, 0, 0, 0, currentLocation)
	endDate = startDate.AddDate(0, 1, -1)

	return
}

func PrevMonths(base time.Time, count int) (startDate time.Time, endDate time.Time) {
	currentLocation := base.Location()

	SearchedMonth := base.AddDate(0, -count, 0)
	startDate = time.Date(SearchedMonth.Year(), SearchedMonth.Month(), 1, 0, 0, 0, 0, currentLocation)
	endDate = startDate.AddDate(0, 1, -1)

	return
}

func NextMonths(base time.Time, count int) (startDate time.Time, endDate time.Time) {
	currentLocation := base.Location()

	SearchedMonth := base.AddDate(0, count, 0)
	startDate = time.Date(SearchedMonth.Year(), SearchedMonth.Month(), 1, 0, 0, 0, 0, currentLocation)
	endDate = startDate.AddDate(0, 1, -1)

	return
}

func ThisYear(base time.Time) (startDate time.Time, endDate time.Time) {
	currentLocation := base.Location()

	startDate = time.Date(base.Year(), 1, 1, 0, 0, 0, 0, currentLocation)
	endDate = startDate.AddDate(1, 0, -1)

	return
}

func PrevYears(base time.Time, count int) (startDate time.Time, endDate time.Time) {
	currentLocation := base.Location()

	SearchedYear := base.AddDate(-count, 0, 0)
	startDate = time.Date(SearchedYear.Year(), 1, 1, 0, 0, 0, 0, currentLocation)
	endDate = startDate.AddDate(1, 0, -1)

	return
}

func NextYears(base time.Time, count int) (startDate time.Time, endDate time.Time) {
	currentLocation := base.Location()

	SearchedYear := base.AddDate(count, 0, 0)
	startDate = time.Date(SearchedYear.Year(), 1, 1, 0, 0, 0, 0, currentLocation)
	endDate = startDate.AddDate(1, 0, -1)

	return
}
