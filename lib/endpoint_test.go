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
	if result := TrimSerial("11_20170402092959.ts"); result != "20170402092959.ts" {
		t.Errorf("TrimSerial %v", result)
	}
}
