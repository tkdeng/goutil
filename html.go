package goutil

import (
	"bytes"

	"github.com/AspieSoft/go-regex-re2/v2"
)

type encodeHtml struct{}

var HTML *encodeHtml = &encodeHtml{}

var regEscHTML *regex.Regexp = regex.Comp(`[<>&]`)
var regEscFixAmp *regex.Regexp = regex.Comp(`&amp;(amp;)*`)

// EscapeHTML replaces HTML characters with html entities
//
// Also prevents and removes &amp;amp; from results
func (encHtml *encodeHtml) Escape(html []byte) []byte {
	html = regEscHTML.RepFunc(html, func(data func(int) []byte) []byte {
		if bytes.Equal(data(0), []byte("<")) {
			return []byte("&lt;")
		} else if bytes.Equal(data(0), []byte(">")) {
			return []byte("&gt;")
		}
		return []byte("&amp;")
	})
	return regEscFixAmp.RepStrLit(html, []byte("&amp;"))
}

var regEscHTMLArgs *regex.Regexp = regex.Comp(`([\\]*)([\\"'\'])`)

// EscapeHTMLArgs escapes quotes and backslashes for use within HTML quotes
// @quote can be used to only escape specific quotes or chars
func (encHtml *encodeHtml) EscapeArgs(html []byte, quote ...byte) []byte {
	if len(quote) == 0 {
		quote = []byte("\"'`")
	}

	return regEscHTMLArgs.RepFunc(html, func(data func(int) []byte) []byte {
		if len(data(1))%2 == 0 && bytes.ContainsRune(quote, rune(data(2)[0])) {
			// return append([]byte("\\"), data(2)...)
			return regex.JoinBytes(data(1), '\\', data(2))
		}
		if bytes.ContainsRune(quote, rune(data(2)[0])) {
			return append(data(1), data(2)...)
		}
		return data(0)
	})
}
