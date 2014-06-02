package riemann

import (
	"encoding/binary"
	"io"
)

var (
	byteOrder = binary.BigEndian
)

type header struct {
	Length uint32
}

func Unpack(reader io.Reader) (data []byte, err error) {
	size, err := UnpackHeader(reader)
	if err != nil {
		return nil, err
	}
	return UnpackData(reader, size)
}

func UnpackData(reader io.Reader, size uint32) (data []byte, err error) {
	data = make([]byte, size)
	if err := readn(reader, data, size); err != nil {
		return nil, err
	}
	return data, nil
}

func UnpackHeader(reader io.Reader) (size uint32, err error) {
	h := header{}
	if err = binary.Read(reader, byteOrder, &h); err != nil {
		return 0, err
	}
	return h.Length, nil
}

func readn(reader io.Reader, buf []byte, len uint32) error {
	rb := int64(len)
	for rb > 0 {
		read, err := reader.Read(buf)
		if err != nil {
			return err
		}
		rb -= int64(read)
		buf = buf[read:]
	}
	return nil
}
