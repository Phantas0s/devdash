package gokit

import "math"

// good old utiles stuff
// expand the go language (ONLY very very general abstractions)

func Round(n float64, precision int) (r float64) {
	pow := math.Pow(10, float64(precision))
	num := n * pow
	_, div := math.Modf(num)

	if n >= 0 && div >= 0.5 {
		r = math.Ceil(num)
	} else if n < 0 && div > -0.5 {
		r = math.Ceil(num)
	} else {
		r = math.Floor(num)
	}

	return r / pow
}

func ConvertBinUnit(val float64, base, to string) (tv float64) {
	convert := map[string]int{
		"b":  6,
		"kb": 5,
		"mb": 4,
		"gb": 3,
		"tb": 2,
		"pb": 1,
		"eb": 0,
	}

	r := convert[base] - convert[to]

	if r >= 0 {
		return Round(val/(math.Pow(float64(1024), float64(r))), 2)
	}

	// TODO to test
	if r < 0 {
		tv = Round(val*(math.Pow(float64(1024), float64(r))), 2)
	}

	return
}
