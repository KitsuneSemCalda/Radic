package types

import "testing"

func TestPinnedWidth(t *testing.T) {
	tests := []struct {
		name string
		want Mold
	}{
		{"i8", i8},
		{"i16", i16},
		{"i32", i32},
		{"i64", i64},
		{"u8", u8},
		{"u16", u16},
		{"u32", u32},
		{"u64", u64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := PinnedWidth(tt.name)
			if !ok {
				t.Fatalf("PinnedWidth(%q) ok = false, want true", tt.name)
			}
			if got != tt.want {
				t.Errorf("PinnedWidth(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPinnedWidthUnknown(t *testing.T) {
	tests := []string{"", "int", "i128", "U8", "float"}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			if _, ok := PinnedWidth(name); ok {
				t.Errorf("PinnedWidth(%q) ok = true, want false", name)
			}
		})
	}
}
