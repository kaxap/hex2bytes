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