package tracer

import (
	"encoding/json"
	"io"

	"github.com/ugorji/go/codec"
)

// encoding specifies a supported encoding type that can be used by encodedReader.
type encoding int

const (
	// encodingMsgpack is "msgpack" encoding.
	encodingMsgpack encoding = iota

	// encodingJSON is JSON encoding.
	encodingJSON
)

// contentType returns the HTTP content type for the encoding.
func (e encoding) contentType() string {
	switch e {
	case encodingMsgpack:
		return "application/msgpack"
	case encodingJSON:
		return "application/json"
	}
	return ""
}

var mh codec.MsgpackHandle

// encodedReader provides an io.Reader which reads from 'v' transformed through the given encoding.
func encodedReader(e encoding, v interface{}) io.Reader {
	r, w := io.Pipe()
	var encoder interface {
		Encode(v interface{}) error
	}
	switch e {
	case encodingJSON:
		encoder = json.NewEncoder(w)
	case encodingMsgpack:
		encoder = codec.NewEncoder(w, &mh)
	default:
		panic("unsupported encoding")
	}
	go func() {
		w.CloseWithError(encoder.Encode(v))
	}()
	return r
}
