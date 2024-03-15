package float

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCastWithDecimal(t *testing.T) {
	cases := []struct {
		Origin      float64
		Decimal     int
		ExpectRound float64
		ExpectFloor float64
	}{
		{1.23456, 2, 1.23, 1.23},
		{1.23546, 2, 1.24, 1.23},
		{1.2, 2, 1.2, 1.2},
		{1, 2, 1, 1},
	}
	for _, v := range cases {
		round := CastWithDecimal(v.Origin, v.Decimal, true)
		assert.Equal(t, v.ExpectRound, round, "round failed")
		floor := CastWithDecimal(v.Origin, v.Decimal, false)
		assert.Equal(t, v.ExpectFloor, floor, "floor failed")
	}
}

// BenchmarkCastWithDecimalNew-8           1000000000               0.000000 ns/op
func BenchmarkCastWithDecimal(b *testing.B) {
	cases := []struct {
		Origin      float64
		Decimal     int
		ExpectRound float64
		ExpectFloor float64
	}{
		{1.23456, 2, 1.23, 1.23},
		{1.23546, 2, 1.24, 1.23},
		{1.2, 2, 1.2, 1.2},
		{1, 2, 1, 1},
	}
	for _, v := range cases {
		_ = CastWithDecimal(v.Origin, v.Decimal, true)
		_ = CastWithDecimal(v.Origin, v.Decimal, false)
	}
}
