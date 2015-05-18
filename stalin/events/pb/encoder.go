package pb

import (
	"bytes"
	"code.google.com/p/gogoprotobuf/proto"
	"encoding/binary"
	"errors"
	"io"
	"sync"
)

type Encoder struct {
	mu  sync.Mutex
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(e interface{}) error {
	msg, ok := e.(*Message)
	if !ok {
		enc.err = errors.New("proto: attempt to encode into wrong type")
		return enc.err
	}
	return enc.EncodeMsg(msg)
}

func (enc *Encoder) EncodeMsg(msg *Message) error {
	enc.mu.Lock()
	defer enc.mu.Unlock()

	buf := &bytes.Buffer{}

	data, err := proto.Marshal(msg)
	if err != nil {
		enc.err = err
		return enc.err
	}

	if enc.encodeHeader(uint32(len(data)), buf); err != nil {
		return enc.err
	}

	if _, enc.err = buf.Write(data); enc.err != nil {
		return enc.err
	}

	_, enc.err = enc.w.Write(buf.Bytes())
	return enc.err
}

func (enc *Encoder) encodeHeader(size uint32, w io.Writer) {
	h := &header{size}
	enc.err = binary.Write(w, binary.BigEndian, h)
}
