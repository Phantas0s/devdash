package humanmath

import "math"

// Do some simple formatting for human readibility
// TODO need to change the title for that

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

// TODO these conversions are pretty ugly
// TODO floating point number (precision 2)
// TODO maybe doing that with bit shifting? -> https://github.com/likexian/gokit/blob/master/xhuman/xhuman.go
func ConvertBinUnit(val []int, base, to string) (tv []int) {
	convert := map[string]int{
		"b":  3,
		"kb": 2,
		"mb": 1,
		"gb": 0,
	}

	r := convert[base] - convert[to]

	if r >= 0 {
		for _, v := range val {
			tv = append(tv, int(Round(float64(v)/(math.Pow(float64(1024), float64(r))), 2)))
		}
	}

	// TODO to test
	if r < 0 {
		for _, v := range val {
			tv = append(tv, int(Round(float64(v)*(math.Pow(float64(1024), float64(r))), 2)))
		}
	}

	return
}
