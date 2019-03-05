package internal

import (
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/totime"
	"github.com/pkg/errors"
)

const (
	today     = "today"
	yesterday = "yesterday"
	days_ago  = "days_ago"

	this_week = "this_week"
	last_week = "last_week"
	weeks_ago = "weeks_ago"

	this_month = "this_month"
	last_month = "last_month"
	months_ago = "months_ago"
)

func ConvertDates(
	base time.Time,
	startDate string,
	endDate string,
) (start time.Time, end time.Time, err error) {
	startDate = resolveAlias(startDate)
	endDate = resolveAlias(endDate)

	start, err = convertStartDate(base, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err = convertEndDate(base, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return
}

func convertStartDate(base time.Time, startDate string) (time.Time, error) {
	if strings.Contains(startDate, today) {
		return base, nil
	}

	if strings.Contains(startDate, days_ago) {
		t := strings.Split(startDate, "_")
		days, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		return base.AddDate(0, 0, -int(days)), nil
	}

	if strings.Contains(startDate, this_week) {
		startDate, _ := totime.ThisWeek(base)
		return startDate, nil
	}

	if strings.Contains(startDate, weeks_ago) {
		t := strings.Split(startDate, "_")
		weeks, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		startDate, _ := totime.PrevWeeks(base, int(weeks))

		return startDate, nil
	}

	if strings.Contains(startDate, this_month) {
		startDate, _ := totime.ThisMonth(base)
		return startDate, nil
	}

	if strings.Contains(startDate, months_ago) {
		t := strings.Split(startDate, "_")
		months, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		startDate, _ := totime.PrevMonths(base, int(months))

		return startDate, nil
	}

	return time.Parse("2006-01-02", startDate)
}

func convertEndDate(base time.Time, endDate string) (time.Time, error) {
	if strings.Contains(endDate, today) {
		return base, nil
	}

	if strings.Contains(endDate, days_ago) {
		t := strings.Split(endDate, "_")
		days, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		return base.AddDate(0, 0, -int(days)), nil
	}

	if strings.Contains(endDate, this_week) {
		_, endDate := totime.ThisWeek(base)
		return endDate, nil
	}

	if strings.Contains(endDate, weeks_ago) {
		t := strings.Split(endDate, "_")
		weeks, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		_, endDate := totime.PrevWeeks(base, int(weeks))

		return endDate, nil
	}

	if strings.Contains(endDate, this_month) {
		_, endDate := totime.ThisMonth(base)
		return endDate, nil
	}

	if strings.Contains(endDate, months_ago) {
		t := strings.Split(endDate, "_")
		months, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		_, endDate := totime.PrevMonths(base, int(months))

		return endDate, nil
	}

	return time.Parse("2006-01-02", endDate)
}

func resolveAlias(date string) string {
	if strings.Contains(date, yesterday) {
		return "1_days_ago"
	}

	if strings.Contains(date, last_week) {
		return "1_weeks_ago"
	}

	if strings.Contains(date, last_month) {
		return "1_months_ago"
	}

	return date
}
