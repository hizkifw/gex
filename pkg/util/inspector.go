package util

import (
	"encoding/binary"
	"math"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

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

	inspectRune := func(rn rune, sz int) string {
		var unicodeStr string
		if rn == utf8.RuneError {
			unicodeStr = "Invalid"
		} else {
			if unicode.IsGraphic(rn) {
				unicodeStr = p.Sprintf("%s (%U, %d bytes)", string(rn), rn, sz)
			} else {
				unicodeStr = p.Sprintf("%U, %d bytes", rn, sz)
			}
		}
		return unicodeStr
	}

	// UTF-8
	rn, sz := utf8.DecodeRune(buf)
	res = append(res, Row{"UTF-8", inspectRune(rn, sz)})

	// UTF-16
	utf16Buf := make([]uint16, 0, len(buf)/2)
	for i := 0; i < len(buf)-1; i += 2 {
		utf16Buf = append(utf16Buf, byteOrder.Uint16(buf[i:]))
	}
	rns := utf16.Decode(utf16Buf)
	if len(rns) > 0 {
		rn = rns[0]
	} else {
		rn = utf8.RuneError
	}
	sz = len(utf16.Encode([]rune{rn})) * 2
	res = append(res, Row{"UTF-16", inspectRune(rn, sz)})

	return res
}
