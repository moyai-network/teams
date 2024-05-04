package glyph

import (
	"strings"
	"sync"
	"unicode"
)

var (
	glyphsMu    sync.Mutex
	glyphsCache = map[string]string{}
)

func glyphFromRune(r rune) (string, bool) {
	switch unicode.ToLower(r) {
	case 'a':
		return "\uE10D", true
	case 'b':
		return "\uE10E", true
	case 'c':
		return "\uE10F", true
	case 'd':
		return "\uE110", true
	case 'e':
		return "\uE111", true
	case 'f':
		return "\uE112", true
	case 'g':
		return "\uE113", true
	case 'h':
		return "\uE114", true
	case 'i':
		return "\uE115", true
	case 'j':
		return "\uE116", true
	case 'k':
		return "\uE117", true
	case 'l':
		return "\uE118", true
	case 'm':
		return "\uE119", true
	case 'n':
		return "\uE11A", true
	case 'o':
		return "\uE11B", true
	case 'p':
		return "\uE11C", true
	case 'q':
		return "\uE11D", true
	case 'r':
		return "\uE11E", true
	case 's':
		return "\uE11F", true
	case 't':
		return "\uE120", true
	case 'u':
		return "\uE121", true
	case 'v':
		return "\uE122", true
	case 'w':
		return "\uE123", true
	case 'x':
		return "\uE124", true
	case 'y':
		return "\uE125", true
	case 'z':
		return "\uE126", true
	default:
		return "", false
	}
}

func Parse(str string) string {
	glyphsMu.Lock()
	cache, ok := glyphsCache[str]
	if ok {
		glyphsMu.Unlock()
		return cache
	}
	glyphsMu.Unlock()

	var s strings.Builder
	for _, c := range str {
		gl, ok := glyphFromRune(c)
		if ok {
			s.WriteString(gl)
		} else {
			s.WriteRune(c)
		}
	}

	glyphsMu.Lock()
	glyphsCache[str] = s.String()
	glyphsMu.Unlock()

	return s.String()
}
