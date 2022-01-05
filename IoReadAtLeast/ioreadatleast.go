package ioreadatleast

import (
	"bytes"
	"io"
)

func IsPNG(r io.Reader) (bool, error) {
	magicnum := []byte{137, 80, 78, 71}
	buf := make([]byte, len(magicnum))

	_, err := io.ReadAtLeast(r, buf, len(buf))
	if err != nil {
		return false, err
	}
	return bytes.Equal(magicnum, buf), nil
}
