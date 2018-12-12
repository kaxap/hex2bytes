[![Build Status](https://travis-ci.org/kaxap/hex2bytes.svg?branch=master)](https://travis-ci.org/kaxap/hex2bytes)
[![Coverage Status](https://coveralls.io/repos/github/kaxap/hex2bytes/badge.svg)](https://coveralls.io/github/kaxap/hex2bytes)

# Info

Decodes hex strings into byte slice. Hex bytes must be delimited with space and be in upper case.

# Usage:

```golang
data := "7F 11 2B 3C"
decoded, err := DecodeSpaceDelimitedHex(data)
if err != nil {
    panic(err)
}

fmt.Println(decoded)
```
