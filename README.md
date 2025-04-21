# Elsys Payload Decoder

A minimal, zero‑dependency Go library to decode Elsys LoRa sensor payloads.  
It exposes a simple `Decoder` API that takes a Base64 string and returns a slice of `Measurement` structs.

---

## Installation

```bash
go get github.com/elchemista/elsys_sensor
```

(Replace `github.com/elchemista/elsys_sensor` with your module path.)

---

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/elchemista/elsys_sensor"
)

func main() {
    // create a new decoder
    d := elsys.NewDecoder()

    // base64 payload from your device
    payload := "AQDXAjI="

    // decode to measurements
    measurements, err := d.Decode(payload)
    if err != nil {
        log.Fatalf("decode error: %v", err)
    }

    // print out each metric
    for _, m := range measurements {
        fmt.Printf("%-12s = %.2f\n", m.Name, m.Value)
    }
}
```

**Example output:**

```text
temperature  = 21.50
humidity     = 50.00
```

---

## API

```go
type Measurement struct {
    Name  string  // metric key, e.g. "temperature", "co2", "vdd"
    Value float64 // value in appropriate units
}

type Decoder struct{}

func NewDecoder() *Decoder
    // returns a fresh Decoder

func (d *Decoder) Decode(b64 string) ([]Measurement, error)
    // parses an Elsys Base64 payload into typed measurements
```

## Supported Sensor Types

| Type Code | Name         | Bytes | Units    |
| --------- | ------------ | ----- | -------- |
| `0x01`    | temperature  | 2     | °C /10   |
| `0x02`    | humidity     | 1     | %        |
| `0x03`    | acceleration | 3     | g (±1 g) |
| `0x04`    | light        | 2     | lux      |
| `0x05`    | motion       | 1     | count    |
| `0x06`    | co2          | 2     | ppm      |
| `0x07`    | vdd          | 2     | mV       |
| `0x14`    | pressure     | 4     | hPa      |
| `0x1B`    | analog_uv    | 4     | μV       |
| `0x1C`    | tvoc         | 2     | ppb      |

_Extend the `switch` in `Decode` to support more types as needed._

---

## License

MIT [LICENSE](LICENSE)

---
