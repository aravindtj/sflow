package sflow

import (
	"os"
	"testing"
)

func TestDecode1(t *testing.T) {
	f, err := os.Open("_test/counter_sample.dump")
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(f)

	dgram, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}

	if dgram.Version != 5 {
		t.Errorf("Expected datagram version %v, got %v", 5, dgram.Version)
	}
}
