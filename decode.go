package main

import (
	"errors"
)

var InvalidDataError = errors.New("cannot decode hex: invalid data")

// fromHexChar converts a hex character into its value and a success flag.
// from encoding/hex, upper case only
func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

// DecodeSpaceDelimitedHex decodes a single space-delimited hex string into byte slice
// only uppercase hex letters are supported, for low case hex it returns InvalidByteError
func DecodeSpaceDelimitedHex(s string) ([]byte, error) {

	// special case with 1 byte
	if len(s) == 2 {

		c1, valid := fromHexChar(s[0])
		if !valid {
			return nil, InvalidDataError
		}

		c2, valid := fromHexChar(s[1])
		if !valid {
			return nil, InvalidDataError
		}

		return []byte{c1<<4 | c2}, nil
	}

	readyForFirst := true
	readyForSecond := false
	var b byte = 0
	result := make([]byte, 0, len(s)/3*2)
	for i := 0; i < len(s); i++ {
		if !readyForFirst && !readyForSecond {
			if s[i] == 0x20 {
				if readyForFirst || readyForSecond {
					return nil, InvalidDataError
				}
				readyForFirst = true
				readyForSecond = false
				continue
			} else {
				return nil, InvalidDataError
			}
		}

		if readyForFirst {
			c, valid := fromHexChar(s[i])
			if !valid {
				return nil, InvalidDataError
			}

			b = c << 4
			readyForFirst = false
			readyForSecond = true
			continue
		}

		if readyForSecond {
			c, valid := fromHexChar(s[i])
			if !valid {
				return nil, InvalidDataError
			}

			b |= c
			result = append(result, b)
			readyForSecond = false
		}

	}

	if readyForSecond || readyForFirst {
		return nil, InvalidDataError
	}

	return result, nil

}
