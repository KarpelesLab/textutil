package textutil

import (
	"bytes"
	"unicode"
)

// based on https://github.com/mitchellh/go-wordwrap/blob/master/wordwrap.go (MIT license)

var (
	CRLF = []byte{'\r', '\n'}
	LF   = []byte{'\n'}
)

const nbsp = '\xa0'

// WrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word. pfx can be set to define a prefix for each new line.
func WrapString(s, pfx string, lim int, linebreak []byte) string {
	buf := &bytes.Buffer{}

	var current int
	var wordBuf, spaceBuf bytes.Buffer
	var wordBufLen, spaceBufLen int // number of runes

	for _, char := range s {
		if unicode.IsSpace(char) && char != nbsp {
			switch char {
			case '\n':
				if wordBufLen != 0 {
					spaceBuf.WriteTo(buf)
					wordBuf.WriteTo(buf)
					wordBuf.Reset()
					wordBufLen = 0
				}
				// linebreak, drop white spaces & continue
				spaceBuf.Reset()
				spaceBufLen = 0
				buf.WriteRune(char)
				current = 0
			default:
				if wordBufLen > 0 {
					// flush now
					current += spaceBufLen + wordBufLen
					spaceBuf.WriteTo(buf)
					spaceBuf.Reset()
					spaceBufLen = 0
					wordBuf.WriteTo(buf)
					wordBuf.Reset()
					wordBufLen = 0
					if current >= lim {
						buf.Write(linebreak)
						buf.WriteString(pfx)
						current = len(pfx)
					}
				}

				if current > 0 {
					spaceBuf.WriteRune(char)
					spaceBufLen += 1
				}
			}
			continue
		}
		wordBuf.WriteRune(char)
		wordBufLen += 1

		if current+wordBufLen+spaceBufLen > lim && wordBufLen < lim {
			buf.Write(linebreak)
			buf.WriteString(pfx)
			current = len(pfx)
			spaceBuf.Reset()
			spaceBufLen = 0
		}
	}

	if wordBufLen > 0 {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}
