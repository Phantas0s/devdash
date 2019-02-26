package totime

import (
	"time"
)

const yyyymmdd = "2006-01-02"

func ThisWeek() (startDate string, endDate string) {
	tm := time.Now()
	s := time.Time{}

	// 1 = Monday
	weekDay := int(tm.Weekday())
	if weekDay > 1 {
		s = tm.AddDate(0, 0, -(weekDay - 1))
	}

	if weekDay == 0 {
		s = tm.AddDate(0, 0, weekDay+1)
	}

	startDate = s.Format(yyyymmdd)
	endDate = s.AddDate(0, 0, 7).Format(yyyymmdd)

	return
}

func NPrevWeek(count int) (startDate string, endDate string) {
	tm := time.Now()
	s := time.Time{}

	// 1 = Monday
	weekDay := int(tm.Weekday())
	if weekDay > 1 {
		s = tm.AddDate(0, 0, -((weekDay - 1) * (-count)))
	}

	if weekDay == 0 {
		s = tm.AddDate(0, 0, (weekDay+1)*(-count))
	}

	startDate = s.Format(yyyymmdd)
	endDate = s.AddDate(0, 0, 7*(-count)).Format(yyyymmdd)

	return
}

func NNextWeek(count int) (startDate string, endDate string) {
	tm := time.Now()
	s := time.Time{}

	// 1 = Monday
	weekDay := int(tm.Weekday())
	if weekDay > 1 {
		s = tm.AddDate(0, 0, -((weekDay - 1) * count))
	}

	if weekDay == 0 {
		s = tm.AddDate(0, 0, (weekDay+1)*count)
	}

	startDate = s.Format(yyyymmdd)
	endDate = s.AddDate(0, 0, 7*count).Format(yyyymmdd)

	return
}

func ThisMonth() (startDate string, endDate string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	startDate = firstOfMonth.Format(yyyymmdd)
	endDate = lastOfMonth.Format(yyyymmdd)

	return
}

func NPrevMonth(count int) (startDate string, endDate string) {
	now := time.Now()
	currentLocation := now.Location()

	SearchedMonth := now.AddDate(0, -count, 0)
	firstOfMonth := time.Date(SearchedMonth.Year(), SearchedMonth.Month(), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	startDate = firstOfMonth.Format(yyyymmdd)
	endDate = lastOfMonth.Format(yyyymmdd)

	return
}

func NNextMonth(count int) (startDate string, endDate string) {
	now := time.Now()
	currentLocation := now.Location()

	SearchedMonth := now.AddDate(0, count, 0)
	firstOfMonth := time.Date(SearchedMonth.Year(), SearchedMonth.Month(), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	startDate = firstOfMonth.Format(yyyymmdd)
	endDate = lastOfMonth.Format(yyyymmdd)

	return
}
