package types

import (
	"fmt"
	"math"
)

// signedMins, signedMaxes, and unsignedMaxs hold, at index i, the
// inclusive bound for width 8<<i (i.e. 8, 16, 32, 64) of a signed
// value's minimum, a signed value's maximum, and an unsigned value's
// maximum respectively.
var (
	signedMins   = [4]int64{-128, -32768, -(int64(1) << 31), math.MinInt64}
	signedMaxes  = [4]uint64{127, 32767, (uint64(1) << 31) - 1, (uint64(1) << 63) - 1}
	unsignedMaxs = [4]uint64{255, 65535, (uint64(1) << 32) - 1, math.MaxUint64}
)

// widthIndex maps a mold's bit width to its index into signedMins,
// signedMaxes, and unsignedMaxs.
func widthIndex(width uint8) int {
	switch width {
	case 8:
		return 0
	case 16:
		return 1
	case 32:
		return 2
	default:
		return 3
	}
}

// rangeMin returns the smallest value m can represent: 0 for an
// unsigned mold, or its signed minimum otherwise.
func rangeMin(m Mold) int64 {
	if !m.Signed {
		return 0
	}
	return signedMins[widthIndex(m.Width)]
}

// rangeMax returns the largest value m can represent.
func rangeMax(m Mold) uint64 {
	if m.Signed {
		return signedMaxes[widthIndex(m.Width)]
	}
	return unsignedMaxs[widthIndex(m.Width)]
}

// CombineMold returns the smallest mold whose range contains the
// ranges of both a and b, e.g. for use when unifying the operand
// molds of a binary expression. If either input is signed, the result
// is signed; a signed mold can only be produced up to i64, so it
// returns an error when the combined range needs a signed value below
// i64's minimum or an unsigned value above i64's maximum.
func CombineMold(a, b Mold) (Mold, error) {
	lo := rangeMin(a)
	if bl := rangeMin(b); bl < lo {
		lo = bl
	}
	hi := rangeMax(a)
	if bh := rangeMax(b); bh > hi {
		hi = bh
	}

	for i, w := range []uint8{8, 16, 32, 64} {
		if lo < 0 {
			if lo >= signedMins[i] && hi <= signedMaxes[i] {
				return Mold{Width: w, Signed: true}, nil
			}
			continue
		}
		if hi <= unsignedMaxs[i] {
			return Mold{Width: w, Signed: false}, nil
		}
	}
	return Mold{}, fmt.Errorf("a signed value and a value above i64's max cannot share one mold")
}
