package util

import (
	"encoding/binary"
	"math"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/unicode/runenames"
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

	inspectRune := func(rn rune, sz int) (string, string) {
		var descr string
		name := runenames.Name(rn)
		s := "s"
		if sz == 1 {
			s = ""
		}
		if unicode.IsGraphic(rn) {
			descr = p.Sprintf("%s (%U, %d byte%s)", string(rn), rn, sz, s)
		} else {
			descr = p.Sprintf("%U, %d byte%s", rn, sz, s)
		}
		return descr, name
	}

	// UTF-8
	rn, sz := utf8.DecodeRune(buf)
	if rn != utf8.RuneError {
		descr, name := inspectRune(rn, sz)
		res = append(res, Row{"UTF-8", descr}, Row{"", name})
	}

	// UTF-16
	utf16Buf := make([]uint16, 0, len(buf)/2)
	for i := 0; i < Min(len(buf)-1, 4); i += 2 {
		utf16Buf = append(utf16Buf, byteOrder.Uint16(buf[i:]))
	}
	rns := utf16.Decode(utf16Buf)
	if len(rns) > 0 {
		rn = rns[0]
	} else {
		rn = utf8.RuneError
	}
	sz = len(utf16.Encode([]rune{rn})) * 2
	if rn != utf8.RuneError {
		descr, name := inspectRune(rn, sz)
		res = append(res, Row{"UTF-16", descr}, Row{"", name})
	}

	return res
}
