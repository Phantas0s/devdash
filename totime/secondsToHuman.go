package totime

// TODO to complete
// TODO propose condensed version 8y3m9d3h2m???
// import (
// 	"math"
// )

// const (
// 	// display
// 	year   = "year"
// 	month  = "month"
// 	week   = "week"
// 	day    = "day"
// 	hour   = "hour"
// 	minute = "minute"
// 	second = "second"
// )

// func display(word string, count int) (result string) {
// 	if count > 1 {
// 		return word + "s"
// 	}
// }

// func secondsToHuman(input int) (result string) {
// 	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
// 	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)

// 	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
// 	seconds = input % (60 * 60 * 24 * 7 * 30)

// 	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
// 	seconds = input % (60 * 60 * 24 * 7)

// 	days := math.Floor(float64(seconds) / 60 / 60 / 24)
// 	seconds = input % (60 * 60 * 24)

// 	hours := math.Floor(float64(seconds) / 60 / 60)
// 	seconds = input % (60 * 60)

// 	minutes := math.Floor(float64(seconds) / 60)
// 	seconds = input % 60

// 	if years > 0 {
// 		result = int(years), display(year, years)
// 	}

// 	if months > 0 {
// 		result += int(months, "month")
// 	}

// 	if weeks > 0 {
// 		result = int(weeks) + "week")
// 	}

// 	if days > 0 {
// 		result = plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
// 	}

// 	if hours > 0 {
// 		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
// 	}

// 	if minutes > 0 {
// 		result = plural(int(minutes), "minute") + plural(int(seconds), "second")
// 	}
// 	k
// 	result = plural(int(seconds), "second")

// 	return
// }
