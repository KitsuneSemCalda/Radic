package types

import "testing"

func TestCombineMold(t *testing.T) {
	tests := []struct {
		name string
		a, b Mold
		want Mold
	}{
		{"same mold", i32, i32, i32},
		{"unsigned widens to fit other unsigned", u8, u16, u16},
		{"signed widens to fit other signed", i8, i32, i32},
		{"unsigned promoted to hold negative signed range", u8, i8, i16},
		{"u32 and i32 need i64 to hold both ranges", u32, i32, i64},
		{"order does not matter", i32, u8, i32},
		{"u64 and u8 stay unsigned", u64, u8, u64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineMold(tt.a, tt.b)
			if err != nil {
				t.Fatalf("CombineMold(%v, %v) returned error: %v", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("CombineMold(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCombineMoldOverflow(t *testing.T) {
	if _, err := CombineMold(u64, i8); err == nil {
		t.Error("CombineMold(u64, i8) returned nil error, want non-nil")
	}
}
