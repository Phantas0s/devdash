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
	daysAgo   = "days_ago"

	thisWeek = "this_week"
	lastWeek = "last_week"
	weeksAgo = "weeks_ago"

	thisMonth = "this_month"
	lastMonth = "last_month"
	monthsAgo = "months_ago"

	thisYear = "this_year"
	lastYear = "last_year"
	yearsAgo = "years_ago"
)

// ConvertDates from configuration string values to formatted start date / end date with layout.
// Example: "next_month" => startDate "2019-01-01", endDate "2019-01-31".
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

	if strings.Contains(startDate, daysAgo) {
		t := strings.Split(startDate, "_")
		days, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		return base.AddDate(0, 0, -int(days)), nil
	}

	if strings.Contains(startDate, thisWeek) {
		startDate, _ := totime.ThisWeek(base)
		return startDate, nil
	}

	if strings.Contains(startDate, weeksAgo) {
		t := strings.Split(startDate, "_")
		weeks, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		startDate, _ := totime.PrevWeeks(base, int(weeks))

		return startDate, nil
	}

	if strings.Contains(startDate, thisMonth) {
		startDate, _ := totime.ThisMonth(base)
		return startDate, nil
	}

	if strings.Contains(startDate, monthsAgo) {
		t := strings.Split(startDate, "_")
		months, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		startDate, _ := totime.PrevMonths(base, int(months))

		return startDate, nil
	}

	if strings.Contains(startDate, thisYear) {
		startDate, _ := totime.ThisYear(base)
		return startDate, nil
	}

	if strings.Contains(startDate, yearsAgo) {
		t := strings.Split(startDate, "_")
		years, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", startDate)
		}

		startDate, _ := totime.PrevYears(base, int(years))

		return startDate, nil
	}

	return time.Parse("2006-01-02", startDate)
}

func convertEndDate(base time.Time, endDate string) (time.Time, error) {
	if strings.Contains(endDate, today) {
		return base, nil
	}

	if strings.Contains(endDate, daysAgo) {
		t := strings.Split(endDate, "_")
		days, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		return base.AddDate(0, 0, -int(days)), nil
	}

	if strings.Contains(endDate, thisWeek) {
		_, endDate := totime.ThisWeek(base)
		return endDate, nil
	}

	if strings.Contains(endDate, weeksAgo) {
		t := strings.Split(endDate, "_")
		weeks, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		_, endDate := totime.PrevWeeks(base, int(weeks))

		return endDate, nil
	}

	if strings.Contains(endDate, thisMonth) {
		_, endDate := totime.ThisMonth(base)
		return endDate, nil
	}

	if strings.Contains(endDate, monthsAgo) {
		t := strings.Split(endDate, "_")
		months, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		_, endDate := totime.PrevMonths(base, int(months))

		return endDate, nil
	}

	if strings.Contains(endDate, thisYear) {
		endDate, _ := totime.ThisYear(base)
		return endDate, nil
	}

	if strings.Contains(endDate, yearsAgo) {
		t := strings.Split(endDate, "_")
		years, err := strconv.ParseInt(t[0], 0, 0)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "%s is not a valid date", endDate)
		}

		endDate, _ := totime.PrevYears(base, int(years))

		return endDate, nil
	}

	return time.Parse("2006-01-02", endDate)
}

func resolveAlias(date string) string {
	if strings.Contains(date, yesterday) {
		return "1_days_ago"
	}

	if strings.Contains(date, lastWeek) {
		return "1_weeks_ago"
	}

	if strings.Contains(date, lastMonth) {
		return "1_months_ago"
	}

	if strings.Contains(date, lastYear) {
		return "1_years_ago"
	}

	return date
}
