package datasize

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ByteSize uint64

const (
	B  ByteSize = 1
	KB          = B << 10
	MB          = KB << 10
	GB          = MB << 10
	TB          = GB << 10
	PB          = TB << 10
	EB          = PB << 10

	fnUnmarshalText string = "UnmarshalText"
	maxUint64       uint64 = (1 << 64) - 1
	cutoff          uint64 = maxUint64 / 10
)

var ErrBits = errors.New("unit with capital unit prefix and lower case unit (b) - bits, not bytes ")

func (b ByteSize) Bytes() uint64 {
	return uint64(b)
}

func (b ByteSize) KBytes() float64 {
	v := b / KB
	r := b % KB
	return float64(v) + float64(r)/float64(KB)
}

func (b ByteSize) MBytes() float64 {
	v := b / MB
	r := b % MB
	return float64(v) + float64(r)/float64(MB)
}

func (b ByteSize) GBytes() float64 {
	v := b / GB
	r := b % GB
	return float64(v) + float64(r)/float64(GB)
}

func (b ByteSize) TBytes() float64 {
	v := b / TB
	r := b % TB
	return float64(v) + float64(r)/float64(TB)
}

func (b ByteSize) PBytes() float64 {
	v := b / PB
	r := b % PB
	return float64(v) + float64(r)/float64(PB)
}

func (b ByteSize) EBytes() float64 {
	v := b / EB
	r := b % EB
	return float64(v) + float64(r)/float64(EB)
}

func (b ByteSize) String() string {
	switch {
	case b == 0:
		return "0B"
	case b%EB == 0:
		return fmt.Sprintf("%dEB", b/EB)
	case b%PB == 0:
		return fmt.Sprintf("%dPB", b/PB)
	case b%TB == 0:
		return fmt.Sprintf("%dTB", b/TB)
	case b%GB == 0:
		return fmt.Sprintf("%dGB", b/GB)
	case b%MB == 0:
		return fmt.Sprintf("%dMB", b/MB)
	case b%KB == 0:
		return fmt.Sprintf("%dKB", b/KB)
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func (b ByteSize) HR() string {
	return b.HumanReadable()
}

func (b ByteSize) HumanReadable() string {
	switch {
	case b > EB:
		return fmt.Sprintf("%.1f EB", b.EBytes())
	case b > PB:
		return fmt.Sprintf("%.1f PB", b.PBytes())
	case b > TB:
		return fmt.Sprintf("%.1f TB", b.TBytes())
	case b > GB:
		return fmt.Sprintf("%.1f GB", b.GBytes())
	case b > MB:
		return fmt.Sprintf("%.1f MB", b.MBytes())
	case b > KB:
		return fmt.Sprintf("%.1f KB", b.KBytes())
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func (b ByteSize) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *ByteSize) UnmarshalText(t []byte) error {
	var val uint64
	var unit string
	var unitUint64 uint64 = 1

	hasDecimal := false
	var decimal uint64
	var power uint64 = 1

	// copy for error message
	t0 := t

	var c byte
	var i int

ParseLoop:
	for i < len(t) {
		c = t[i]
		switch {
		case '0' <= c && c <= '9':
			if !hasDecimal {
				if val > cutoff {
					goto Overflow
				}

				c = c - '0'
				val *= 10

				if val > val+uint64(c) {
					// val+v overflows
					goto Overflow
				}
				val += uint64(c)
			} else {
				if decimal > cutoff {
					goto Overflow
				}

				c = c - '0'
				decimal *= 10

				if decimal > decimal+uint64(c) {
					// decimal+v overflows
					goto Overflow
				}
				decimal += uint64(c)
				power *= 10
			}
			i++
		case c == '.':
			if hasDecimal {
				goto SyntaxError
			}
			hasDecimal = true
			i++
		default:
			if i == 0 {
				goto SyntaxError
			}
			break ParseLoop
		}
	}

	unit = strings.TrimSpace(string(t[i:]))
	switch unit {
	case "Kb", "Mb", "Gb", "Tb", "Pb", "Eb":
		goto BitsError
	}
	unit = strings.ToLower(unit)
	switch unit {
	case "", "b", "byte", "bytes":
		// do nothing - already in bytes

	case "k", "kb", "kib", "kilo", "kilobyte", "kilobytes":
		if val > maxUint64/uint64(KB) {
			goto Overflow
		}
		unitUint64 = uint64(KB)
		val *= unitUint64

	case "m", "mb", "mib", "mega", "megabyte", "megabytes":
		if val > maxUint64/uint64(MB) {
			goto Overflow
		}
		unitUint64 = uint64(MB)
		val *= unitUint64

	case "g", "gb", "gib", "giga", "gigabyte", "gigabytes":
		if val > maxUint64/uint64(GB) {
			goto Overflow
		}
		unitUint64 = uint64(GB)
		val *= unitUint64

	case "t", "tb", "tib", "tera", "terabyte", "terabytes":
		if val > maxUint64/uint64(TB) {
			goto Overflow
		}
		unitUint64 = uint64(TB)
		val *= unitUint64

	case "p", "pb", "pib", "peta", "petabyte", "petabytes":
		if val > maxUint64/uint64(PB) {
			goto Overflow
		}
		unitUint64 = uint64(PB)
		val *= unitUint64

	case "e", "eb", "eib", "exa", "exabyte", "exabytes":
		if val > maxUint64/uint64(EB) {
			goto Overflow
		}
		unitUint64 = uint64(EB)
		val *= unitUint64

	default:
		goto SyntaxError
	}

	decimal = uint64(float64(decimal*unitUint64) / float64(power))
	if decimal > maxUint64/unitUint64 {
		goto Overflow
	}

	*b = ByteSize(val + decimal)
	return nil

Overflow:
	*b = ByteSize(maxUint64)
	return &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: strconv.ErrRange}

SyntaxError:
	*b = 0
	return &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: strconv.ErrSyntax}

BitsError:
	*b = 0
	return &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: ErrBits}
}

func Parse(t []byte) (ByteSize, error) {
	var v ByteSize
	err := v.UnmarshalText(t)
	return v, err
}

func MustParse(t []byte) ByteSize {
	v, err := Parse(t)
	if err != nil {
		panic(err)
	}
	return v
}

func ParseString(s string) (ByteSize, error) {
	return Parse([]byte(s))
}

func MustParseString(s string) ByteSize {
	return MustParse([]byte(s))
}
