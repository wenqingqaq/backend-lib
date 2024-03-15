package float

import (
	"math"
)

func CastWithDecimal(origin float64, decimal int, round bool) float64 {
	multiple := math.Pow10(decimal)
	if round {
		return math.Round(origin*multiple) / multiple
	}
	return math.Floor(origin*multiple) / multiple
}
