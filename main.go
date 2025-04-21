package main

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// Measurement is a named metric value returned by Decode.
type Measurement struct {
	Name  string
	Value float64
}

// Decoder converts a base‑64 payload into a slice of measurements.
type Decoder struct{}

// NewDecoder returns a fresh Decoder.
func NewDecoder() *Decoder { return &Decoder{} }

// Decode parses Elsys TLV data carried in a base‑64 string.
func (d *Decoder) Decode(b64 string) ([]Measurement, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}
	if len(raw) == 0 {
		return nil, errors.New("empty payload")
	}

	var m []Measurement
	for i := 0; i < len(raw); i++ {
		switch raw[i] {

		case TYPE_TEMP: // 2 B, °C /10
			if !have(i, 2, raw) {
				return nil, errShort
			}
			v := int16(raw[i+1])<<8 | int16(raw[i+2])
			m = append(m, Measurement{Name: "temperature", Value: float64(v) / 10})
			i += 2

		case TYPE_RH: // 1 B, %
			if !have(i, 1, raw) {
				return nil, errShort
			}
			m = append(m, Measurement{Name: "humidity", Value: float64(raw[i+1])})
			i++

		case TYPE_ACC: // 3 B, ±g (63 LSB = 1 g)
			if !have(i, 3, raw) {
				return nil, errShort
			}
			x, y, z := int8(raw[i+1]), int8(raw[i+2]), int8(raw[i+3])
			m = append(m,
				Measurement{Name: "acc_x", Value: float64(x) / 63},
				Measurement{Name: "acc_y", Value: float64(y) / 63},
				Measurement{Name: "acc_z", Value: float64(z) / 63},
			)
			i += 3

		case TYPE_LIGHT: // 2 B, lx
			if !have(i, 2, raw) {
				return nil, errShort
			}
			v := u16(raw[i+1], raw[i+2])
			m = append(m, Measurement{Name: "light", Value: float64(v)})
			i += 2

		case TYPE_MOTION: // 1 B, counts
			if !have(i, 1, raw) {
				return nil, errShort
			}
			m = append(m, Measurement{Name: "motion", Value: float64(raw[i+1])})
			i++

		case TYPE_CO2: // 2 B, ppm
			if !have(i, 2, raw) {
				return nil, errShort
			}
			v := u16(raw[i+1], raw[i+2])
			m = append(m, Measurement{Name: "co2", Value: float64(v)})
			i += 2

		case TYPE_VDD: // 2 B, mV
			if !have(i, 2, raw) {
				return nil, errShort
			}
			v := u16(raw[i+1], raw[i+2])
			m = append(m, Measurement{Name: "vdd", Value: float64(v)})
			i += 2

		case TYPE_EXT_ANALOG_UV: // 4 B, µV (signed)
			if !have(i, 4, raw) {
				return nil, errShort
			}
			v := i32(raw[i+1], raw[i+2], raw[i+3], raw[i+4])
			m = append(m, Measurement{Name: "analog_uv", Value: float64(v)})
			i += 4

		case TYPE_TVOC: // 2 B, ppb
			if !have(i, 2, raw) {
				return nil, errShort
			}
			v := u16(raw[i+1], raw[i+2])
			m = append(m, Measurement{Name: "tvoc", Value: float64(v)})
			i += 2

		case TYPE_PRESSURE: // 4 B, hPa * 1000
			if !have(i, 4, raw) {
				return nil, errShort
			}
			v := u32(raw[i+1], raw[i+2], raw[i+3], raw[i+4])
			m = append(m, Measurement{Name: "pressure", Value: float64(v) / 1000})
			i += 4

		default: // skip unknown type‑code, break loop to avoid infinite walk
			return nil, fmt.Errorf("%w: 0x%02X at index %d", errUnknown, raw[i], i)
		}
	}

	return m, nil
}

const (
	TYPE_TEMP          = 0x01
	TYPE_RH            = 0x02
	TYPE_ACC           = 0x03
	TYPE_LIGHT         = 0x04
	TYPE_MOTION        = 0x05
	TYPE_CO2           = 0x06
	TYPE_VDD           = 0x07
	TYPE_EXT_ANALOG_UV = 0x1B
	TYPE_TVOC          = 0x1C
	TYPE_PRESSURE      = 0x14
)

var (
	errShort   = errors.New("payload truncated while decoding field")
	errUnknown = errors.New("unknown sensor type")
)

// u16 returns an unsigned big‑endian uint16 from two bytes.
func u16(h, l byte) uint16 { return uint16(h)<<8 | uint16(l) }

// u32 returns an unsigned big‑endian uint32 from four bytes.
func u32(b3, b2, b1, b0 byte) uint32 {
	return uint32(b3)<<24 | uint32(b2)<<16 | uint32(b1)<<8 | uint32(b0)
}

// i32 returns a signed big‑endian int32 from four bytes.
func i32(b3, b2, b1, b0 byte) int32 { return int32(u32(b3, b2, b1, b0)) }

// have verifies that n more bytes (after the type byte) are available.
func have(idx, needed int, buf []byte) bool { return idx+needed < len(buf) }
