// Package types implements Radic's fixed-width integer type system: the
// Mold representation itself, inference of a literal's mold from its
// text, lookup of an explicitly pinned mold by name, and combination of
// two molds into the smallest one that can represent both.
package types

import "fmt"

// Mold is a fixed-width integer type: a bit width (8, 16, 32, or 64)
// together with its signedness.
type Mold struct {
	Width  uint8
	Signed bool
}

// Bits returns the width of m in bits.
func (m Mold) Bits() int { return int(m.Width) }

// String returns m in Radic's type-name notation, e.g. "i32" or "u8".
func (m Mold) String() string {
	kind := "u"
	if m.Signed {
		kind = "i"
	}

	return fmt.Sprintf("%s%d", kind, m.Width)
}

// The eight molds of the Radic type system, shared as the return values
// of MoldOfLiteral, PinnedWidth, and CombineMold.
var (
	i8  = Mold{Width: 8, Signed: true}
	i16 = Mold{Width: 16, Signed: true}
	i32 = Mold{Width: 32, Signed: true}
	i64 = Mold{Width: 64, Signed: true}
	u8  = Mold{Width: 8, Signed: false}
	u16 = Mold{Width: 16, Signed: false}
	u32 = Mold{Width: 32, Signed: false}
	u64 = Mold{Width: 64, Signed: false}
)
