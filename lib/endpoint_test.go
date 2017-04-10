package lib

import (
	"testing"
)

func TestEndpoint(t *testing.T) {
	e := NewEndpoint("v-low-tokyo1", "2865")
	if uri := e.List(); uri != "https://smartcast.hs.llnwd.net/v-low-tokyo1/2865/2865.txt" {
		t.Errorf("Endpoint.List %v", uri)
	}
	if uri := e.TS("11_20170402092959.ts"); uri != "https://smartcast.hs.llnwd.net/v-low-tokyo1/2865/11_20170402092959.ts" {
		t.Errorf("Endpoint.TS %v", uri)
	}
}

func TestTrimSerial(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "11_20170402092959.ts",
			expected: "20170402092959.ts",
		},
		{
			input:    "no_serial.ts",
			expected: "no_serial.ts",
		},
		{
			input:    "123_long_serial.ts",
			expected: "123_long_serial.ts",
		},
		{
			input:    "1_short_serial.ts",
			expected: "1_short_serial.ts",
		},
	}
	for _, c := range cases {
		if result := TrimSerial(c.input); result != c.expected {
			t.Errorf("TrimSerial %v %v", result, c.expected)
		}
	}
}
