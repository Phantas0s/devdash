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
