package main

import (
	"reflect"
	"testing"
)

func TestDecodeSingleMeasurements(t *testing.T) {
	decoder := NewDecoder()

	tests := []struct {
		name string
		b64  string
		want []Measurement
	}{
		{
			name: "temperature",
			b64:  "AQDX", // TYPE_TEMP + 0x00D7 -> 21.5°C
			want: []Measurement{{Name: "temperature", Value: 21.5}},
		},
		{
			name: "humidity",
			b64:  "AjI=", // TYPE_RH + 50 -> 50%
			want: []Measurement{{Name: "humidity", Value: 50}},
		},
		{
			name: "light",
			b64:  "BAoA", // TYPE_LIGHT + 0x0280 -> 2560 lx
			want: []Measurement{{Name: "light", Value: 2560}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decoder.Decode(tt.b64)
			if err != nil {
				t.Fatalf("Decode(%q) returned error: %v", tt.b64, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decode(%q) = %v, want %v", tt.b64, got, tt.want)
			}
		})
	}
}

func TestDecodeMultipleMeasurements(t *testing.T) {
	// payload: temp=21.5°C, humidity=50%
	// raw: [TYPE_TEMP,0x00,0xD7, TYPE_RH,50]
	b64 := "AQDXAjI="
	decoder := NewDecoder()
	got, err := decoder.Decode(b64)
	if err != nil {
		t.Errorf("decode(%q) error: %v", b64, err)
	}

	want := []Measurement{
		{Name: "temperature", Value: 21.5},
		{Name: "humidity", Value: 50},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("decode(%q) = %v, want %v", b64, got, want)
	}
}
