package iomultireader

import (
	"bytes"
	"errors"
	"io"
)

func NewPNG(r io.Reader) (io.Reader, error) {
	magicnum := []byte{137, 80, 78, 71}
	buf := make([]byte, len(magicnum))

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	if !bytes.Equal(magicnum, buf) {
		return nil, errors.New("not png")
	}
	pngimg := io.MultiReader(bytes.NewReader(magicnum), r)
	return pngimg, nil
}
