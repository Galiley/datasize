package datasize

import (
	"testing"
)

func TestMarshalText(t *testing.T) {
	table := []struct {
		in  ByteSize
		out string
	}{
		{0, "0B"},
		{B, "1B"},
		{KB, "1KB"},
		{MB, "1MB"},
		{GB, "1GB"},
		{TB, "1TB"},
		{PB, "1PB"},
		{EB, "1EB"},
		{400 * TB, "400TB"},
		{2048 * MB, "2GB"},
		{B + KB, "1025B"},
		{MB + 20*KB, "1044KB"},
		{100*MB + KB, "102401KB"},
	}

	for _, tt := range table {
		b, _ := tt.in.MarshalText()
		s := string(b)

		if s != tt.out {
			t.Errorf("MarshalText(%d) => %s, want %s", tt.in, s, tt.out)
		}
	}
}

func TestUnmarshalText(t *testing.T) {
	table := []struct {
		in  string
		err bool
		out ByteSize
	}{
		{"0", false, ByteSize(0)},
		{"0B", false, ByteSize(0)},
		{"0 KB", false, ByteSize(0)},
		{"1", false, B},
		{"1K", false, KB},
		{"2MB", false, 2 * MB},
		{"5 GB", false, 5 * GB},
		{"20480 G", false, 20 * TB},
		{"50 eB", true, ByteSize((1 << 64) - 1)},
		{"200000 pb", true, ByteSize((1 << 64) - 1)},
		{"10 Mb", true, ByteSize(0)},
		{"g", true, ByteSize(0)},
		{"10 kB ", false, 10 * KB},
		{"10 kBs ", true, ByteSize(0)},
		{"1.1.1.1 KB", true, ByteSize(0)},
		{"10.5 MB", false, ByteSize(uint64(10*MB) + uint64(float64(MB/2)))},
	}

	for _, tt := range table {
		t.Run("UnmarshalText "+tt.in, func(t *testing.T) {
			var s ByteSize
			err := s.UnmarshalText([]byte(tt.in))

			if (err != nil) != tt.err {
				t.Errorf("UnmarshalText(%s) => %v, want no error", tt.in, err)
			}

			if s != tt.out {
				t.Errorf("UnmarshalText(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("Parse "+tt.in, func(t *testing.T) {
			s, err := Parse([]byte(tt.in))

			if (err != nil) != tt.err {
				t.Errorf("Parse(%s) => %v, want no error", tt.in, err)
			}

			if s != tt.out {
				t.Errorf("Parse(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("MustParse "+tt.in, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.err {
					t.Errorf("MustParse(%s) => %v, want no error", tt.in, err)
				}
			}()

			s := MustParse([]byte(tt.in))
			if s != tt.out {
				t.Errorf("MustParse(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("ParseString "+tt.in, func(t *testing.T) {
			s, err := ParseString(tt.in)

			if (err != nil) != tt.err {
				t.Errorf("ParseString(%s) => %v, want no error", tt.in, err)
			}

			if s != tt.out {
				t.Errorf("ParseString(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("MustParseString "+tt.in, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.err {
					t.Errorf("MustParseString(%s) => %v, want no error", tt.in, err)
				}
			}()

			s := MustParseString(tt.in)
			if s != tt.out {
				t.Errorf("MustParseString(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
	}
}
