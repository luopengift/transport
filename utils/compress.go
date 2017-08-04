package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

func Compress(p []byte) ([]byte, error) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	if _, err := w.Write(p); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func Decompress(p []byte) ([]byte, error) {
	var out bytes.Buffer
	in := bytes.NewBuffer(p)
	r, err := zlib.NewReader(in)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
