package types

import "testing"

func TestMoldOfLiteralPositive(t *testing.T) {
	tests := []struct {
		text string
		want Mold
	}{
		{"0", i8},
		{"127", i8},
		{"128", u8},
		{"255", u8},
		{"256", i16},
		{"32767", i16},
		{"32768", u16},
		{"65535", u16},
		{"65536", i32},
		{"2147483647", i32},
		{"2147483648", u32},
		{"4294967295", u32},
		{"4294967296", i64},
		{"9223372036854775807", i64},
		{"9223372036854775808", u64},
		{"18446744073709551615", u64},
		{"0xff", u8},
		{"0b1010", i8},
		{"0o17", i8},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got, err := MoldOfLiteral(tt.text, false)
			if err != nil {
				t.Fatalf("MoldOfLiteral(%q, false) returned error: %v", tt.text, err)
			}
			if got != tt.want {
				t.Errorf("MoldOfLiteral(%q, false) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestMoldOfLiteralNegated(t *testing.T) {
	tests := []struct {
		text string
		want Mold
	}{
		{"0", i8},
		{"128", i8},
		{"129", i16},
		{"32768", i16},
		{"32769", i32},
		{"2147483648", i32},
		{"2147483649", i64},
		{"9223372036854775808", i64},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got, err := MoldOfLiteral(tt.text, true)
			if err != nil {
				t.Fatalf("MoldOfLiteral(%q, true) returned error: %v", tt.text, err)
			}
			if got != tt.want {
				t.Errorf("MoldOfLiteral(%q, true) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestMoldOfLiteralNegatedOverflow(t *testing.T) {
	if _, err := MoldOfLiteral("9223372036854775809", true); err == nil {
		t.Error("MoldOfLiteral(\"9223372036854775809\", true) returned nil error, want non-nil")
	}
}

func TestMoldOfLiteralInvalid(t *testing.T) {
	tests := []string{"", "abc", "-1", "18446744073709551616"}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			if _, err := MoldOfLiteral(text, false); err == nil {
				t.Errorf("MoldOfLiteral(%q, false) returned nil error, want non-nil", text)
			}
		})
	}
}
