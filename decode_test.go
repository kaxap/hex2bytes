package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/profile"
	"strings"
	"testing"
)

func byteSliceRepr(b []byte) string {
	result := make([]string, 0, len(b))
	for _, e := range b {
		result = append(result, fmt.Sprintf("%x", e))
	}
	return strings.Join(result, ",")
}

func getValidTestData() map[string][]byte {
	return map[string][]byte{
		"91 77 70 10 00 00 F5": {0x91, 0x77, 0x70, 0x10, 0x00, 0x00, 0xf5},
		"13":                   {0x13},
		"00":                   {0x00},
		"00 00":                {0x00, 0x00},
		"00 00 00":             {0x00, 0x00, 0x00},
		"00 00 00 00":          {0x00, 0x00, 0x00, 0x00},
		"7F D6":                {0x7F, 0xD6},
		"91 77 80 04 21 64 F3 91 77 80 04 21 64 F3 91 77 80 04 21 64 F3 91 77 80 04 21 64 F3": {0x91, 0x77, 0x80, 0x04, 0x21, 0x64, 0xf3, 0x91, 0x77, 0x80, 0x04, 0x21, 0x64, 0xf3,
			0x91, 0x77, 0x80, 0x04, 0x21, 0x64, 0xf3, 0x91, 0x77, 0x80, 0x04, 0x21, 0x64, 0xf3},
	}
}

func getInvalidTestData() map[string]error {
	return map[string]error{
		"91 77 70 10 00 00 F5 ": InvalidDataError,
		"1o":                    InvalidDataError,
		"-0":                    InvalidDataError,
		" 00":                   InvalidDataError,
		"00 ":                   InvalidDataError,
		"ff 00":                 InvalidDataError,
		" FF 00 ":               InvalidDataError,
		"FF00":                  InvalidDataError,
		"ff 00 gg":              InvalidDataError,
		"91 77 80 04 21 64 F3 91 77 80 04 21 64 F3 91 77 80 04 21 64 F3 91 77 80 04 21 64 F": InvalidDataError,
	}
}

func TestDecodeSpaceDelimitedHex2(t *testing.T) {
	data := "7F D6"
	expected := []byte{0x7f, 0xd6}
	decoded, err := DecodeSpaceDelimitedHex(data)

	if err != nil {
		t.Errorf("Got unexpected error while parsing \"%s\": %s\n", data, err)
		return
	}

	if bytes.Compare(decoded, expected) != 0 {
		t.Errorf("Could not properly parse \"%s\": Expected [%s], got [%s]\n",
			data, byteSliceRepr(expected), byteSliceRepr(decoded))
	}

}

func TestDecodeSpaceDelimitedHexLeadingSpaceErr(t *testing.T) {
	data := " 00"
	_, err := DecodeSpaceDelimitedHex(data)

	if err == nil {
		t.Errorf("Expected error while parsing \"%s\": but got success\n", data)
		return
	}

	if err != InvalidDataError {
		t.Errorf("Expected error while parsing \"%s\": got %s\n", data, err)
		return
	}

}

func TestDecodeSpaceDelimitedHexTrailingSpaceErr(t *testing.T) {
	data := "00 "
	_, err := DecodeSpaceDelimitedHex(data)

	if err == nil {
		t.Errorf("Expected error while parsing \"%s\": but got success\n", data)
		return
	}

	if err != InvalidDataError {
		t.Errorf("Expected error while parsing \"%s\": but got %s\n", data, err)
		return
	}

}

func TestDecodeSpaceDelimitedHexNoSpaceErr(t *testing.T) {
	data := "FF00"
	_, err := DecodeSpaceDelimitedHex(data)

	if err == nil {
		t.Errorf("Expected error while parsing \"%s\": but got success\n", data)
		return
	}

	if err != InvalidDataError {
		t.Errorf("Expected error while parsing \"%s\": but got %s\n", data, err)
		return
	}

}

func TestDecodeSpaceDelimitedHex(t *testing.T) {
	valid := getValidTestData()
	invalid := getInvalidTestData()

	for data, expected := range valid {
		decoded, err := DecodeSpaceDelimitedHex(data)

		if err != nil {
			t.Errorf("Got unexpected error while parsing \"%s\": %s\n", data, err)
			continue
		}

		if bytes.Compare(decoded, expected) != 0 {
			t.Errorf("Could not properly parse \"%s\": Expected [%s], got [%s]\n",
				data, byteSliceRepr(expected), byteSliceRepr(decoded))
			continue
		}

	}

	for data, result := range invalid {
		_, err := DecodeSpaceDelimitedHex(data)

		if err == nil {
			t.Errorf("Invalid data returned no error: \"%s\"\n", data)
			continue
		}

		if err != result {
			t.Errorf("Got unexpected error for invalid data: %s, expected %s for data %s",
				err, result, data)
		}

	}
}

func TestDecodeSpaceDelimitedHexTwoBytes(t *testing.T) {

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			data := strings.ToUpper(fmt.Sprintf("%.2x %.2x", i, j))
			b, err := DecodeSpaceDelimitedHex(data)

			if err != nil {
				t.Errorf("Could not decode data %s, error: %s", data, err)
				return
			}

			if int(b[0]) != i || int(b[1]) != j {
				t.Errorf("Error parsing data: expected [%.2x, %.2x], got [%.2x, %.2x]",
					i, j, b[0], b[1])
				return
			}
		}
	}
}

func TestDecodeSpaceDelimitedHexThreeBytes(t *testing.T) {

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			for k := 0; k < 256; k++ {
				data := strings.ToUpper(fmt.Sprintf("%.2x %.2x %.2x", i, j, k))
				b, err := DecodeSpaceDelimitedHex(data)
				if err != nil {
					t.Errorf("Could not decode data %s, error: %s", data, err)
					return
				}

				if int(b[0]) != i || int(b[1]) != j || int(b[2]) != k {
					t.Errorf("Error parsing data: expected [%.2x, %.2x, %.2x], got [%.2x, %.2x, %.2x]",
						i, j, k, b[0], b[1], b[2])
					return
				}

			}
		}
	}
}

func TestDecodeSpaceDelimitedHexFourBytes(t *testing.T) {

	defer profile.Start(profile.MemProfile).Stop()
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			for k := 0; k < 256; k++ {
				for l := 0; l < 256; l++ {

					data := strings.ToUpper(fmt.Sprintf("%.2x %.2x %.2x %.2x", i, j, k, l))
					b, err := DecodeSpaceDelimitedHex(data)
					if err != nil {
						t.Errorf("Could not decode data %s, error: %s", data, err)
						return
					}

					if int(b[0]) != i || int(b[1]) != j || int(b[2]) != k || int(b[3]) != l {
						t.Errorf("Error parsing data: expected [%.2x, %.2x, %.2x, %.2x], got [%.2x, %.2x, %.2x, %.2x]",
							i, j, k, l, b[0], b[1], b[2], b[3])
						return
					}

				}
			}
		}
	}
}
