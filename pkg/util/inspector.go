package util

import (
	"encoding/binary"
	"math"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Row struct {
	Key string
	Val string
}

// Inspect returns a list of string representations of the byte slice.
func Inspect(buf []byte, byteOrder binary.ByteOrder) []Row {
	p := message.NewPrinter(language.English)
	res := make([]Row, 0)

	if len(buf) >= 1 {
		res = append(res, Row{"Binary", p.Sprintf("%08b", buf[0])})
		res = append(res, Row{"Uint8", p.Sprintf("%d", buf[0])})
		res = append(res, Row{"Int8", p.Sprintf("%d", int8(buf[0]))})
	}

	if len(buf) >= 2 {
		res = append(res, Row{"Uint16", p.Sprintf("%d", byteOrder.Uint16(buf))})
		res = append(res, Row{"Int16", p.Sprintf("%d", int16(byteOrder.Uint16(buf)))})
	}

	if len(buf) >= 4 {
		res = append(res, Row{"Uint32", p.Sprintf("%d", byteOrder.Uint32(buf))})
		res = append(res, Row{"Int32", p.Sprintf("%d", int32(byteOrder.Uint32(buf)))})
		res = append(res, Row{"Float32", p.Sprintf("%f", math.Float32frombits(byteOrder.Uint32(buf)))})
	}

	if len(buf) >= 8 {
		res = append(res, Row{"Uint64", p.Sprintf("%d", byteOrder.Uint64(buf))})
		res = append(res, Row{"Int64", p.Sprintf("%d", int64(byteOrder.Uint64(buf)))})
		res = append(res, Row{"Float64", p.Sprintf("%f", math.Float64frombits(byteOrder.Uint64(buf)))})
	}

	return res
}
