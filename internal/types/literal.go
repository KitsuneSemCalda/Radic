package types

import (
	"fmt"
	"strconv"
	"strings"
)

// MoldOfLiteral infers the smallest mold that can hold the integer
// literal text (decimal, or 0x/0b/0o prefixed), given whether it is
// negated by a preceding unary minus. For a negated literal, magnitude
// is the value's absolute value and the result is always signed; e.g.
// magnitude 128 negated yields i8 since it fits -128. For a
// non-negated literal, the result is the smallest of i8, u8, i16, u16,
// i32, u32, i64, u64 that can hold the value, preferring signed molds
// at each width. It returns an error if text cannot be parsed as an
// unsigned 64-bit integer, or if a negated literal's magnitude exceeds
// what i64 can represent.
func MoldOfLiteral(text string, negated bool) (Mold, error) {
	magnitude, err := parseMagnitude(text)
	if err != nil {
		return Mold{}, err
	}
	if negated {
		switch {
		case magnitude <= 128:
			return i8, nil
		case magnitude <= 32768:
			return i16, nil
		case magnitude <= uint64(1)<<31:
			return i32, nil
		case magnitude <= uint64(1)<<63:
			return i64, nil
		default:
			return Mold{}, fmt.Errorf("literal %q is below the smallest signed fit (i64)", text)
		}
	}
	switch {
	case magnitude <= 127:
		return i8, nil
	case magnitude <= 255:
		return u8, nil
	case magnitude <= 32767:
		return i16, nil
	case magnitude <= 65535:
		return u16, nil
	case magnitude <= (uint64(1)<<31)-1:
		return i32, nil
	case magnitude <= (uint64(1)<<32)-1:
		return u32, nil
	case magnitude <= (uint64(1)<<63)-1:
		return i64, nil
	default:
		return u64, nil
	}
}

// parseMagnitude parses text as an unsigned 64-bit integer, recognizing
// the 0x, 0b, and 0o prefixes for hexadecimal, binary, and octal
// literals respectively and defaulting to decimal otherwise.
func parseMagnitude(text string) (uint64, error) {
	body, base := text, 10
	switch {
	case strings.HasPrefix(text, "0x"), strings.HasPrefix(text, "0X"):
		body, base = text[2:], 16
	case strings.HasPrefix(text, "0b"), strings.HasPrefix(text, "0B"):
		body, base = text[2:], 2
	case strings.HasPrefix(text, "0o"), strings.HasPrefix(text, "0O"):
		body, base = text[2:], 8
	}
	value, err := strconv.ParseUint(body, base, 64)
	if err != nil {
		return 0, fmt.Errorf("literal %q: %v", text, err)
	}
	return value, nil
}
