package types

// PinnedWidth looks up the mold explicitly named by an i8/i16/.../u64
// type annotation. It returns the zero Mold and false if name does not
// name one of Radic's eight fixed-width integer types.
func PinnedWidth(name string) (Mold, bool) {
	switch name {
	case "i8":
		return i8, true
	case "i16":
		return i16, true
	case "i32":
		return i32, true
	case "i64":
		return i64, true
	case "u8":
		return u8, true
	case "u16":
		return u16, true
	case "u32":
		return u32, true
	case "u64":
		return u64, true
	}
	return Mold{}, false
}
