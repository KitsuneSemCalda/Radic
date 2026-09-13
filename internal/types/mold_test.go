package types

import "testing"

func TestMoldBits(t *testing.T) {
	tests := []struct {
		mold Mold
		want int
	}{
		{i8, 8},
		{u8, 8},
		{i16, 16},
		{u16, 16},
		{i32, 32},
		{u32, 32},
		{i64, 64},
		{u64, 64},
	}

	for _, tt := range tests {
		t.Run(tt.mold.String(), func(t *testing.T) {
			if got := tt.mold.Bits(); got != tt.want {
				t.Errorf("Bits() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMoldString(t *testing.T) {
	tests := []struct {
		mold Mold
		want string
	}{
		{i8, "i8"},
		{i16, "i16"},
		{i32, "i32"},
		{i64, "i64"},
		{u8, "u8"},
		{u16, "u16"},
		{u32, "u32"},
		{u64, "u64"},
		{Mold{Width: 32, Signed: false}, "u32"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mold.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
