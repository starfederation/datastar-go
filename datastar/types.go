package datastar

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	NewLine       = "\n"
	DoubleNewLine = "\n\n"
)

var (
	newLineBuf       = []byte(NewLine)
	doubleNewLineBuf = []byte(DoubleNewLine)
)

// The SSE spec ends a line at CRLF, LF, or a bare CR. Content split into data
// lines must be split on all three, or a bare CR lets the client see a line
// the SDK never prefixed.
var lineTerminatorsToNewLine = strings.NewReplacer("\r\n", NewLine, "\r", NewLine)

func normalizeLineTerminators(b []byte) []byte {
	if bytes.IndexByte(b, '\r') < 0 {
		return b
	}
	return []byte(lineTerminatorsToNewLine.Replace(string(b)))
}

func rejectLineTerminators(field, value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("%s must not contain a carriage return or line feed: %q", field, value)
	}
	return nil
}

type flusher interface {
	Flush() error
}
