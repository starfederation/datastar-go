package datastar

import (
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// sseEndOfLine matches every line terminator the SSE spec allows: CRLF, LF, or CR.
var sseEndOfLine = regexp.MustCompile(`\r\n|\r|\n`)

func newRecordedSSE(t *testing.T) (*ServerSentEventGenerator, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	return NewSSE(w, httptest.NewRequest("GET", "/", nil)), w
}

// clientEvents splits the stream the way an EventSource parser does and
// returns the field lines of each dispatched event.
func clientEvents(stream string) [][]string {
	var events [][]string
	var current []string
	for _, line := range sseEndOfLine.Split(stream, -1) {
		if line == "" {
			if current != nil {
				events = append(events, current)
				current = nil
			}
			continue
		}
		current = append(current, line)
	}
	return events
}

func TestPatchElementsSplitsOnEverySSELineTerminator(t *testing.T) {
	sse, w := newRecordedSSE(t)

	if err := sse.PatchElements("<p>a\rb\r\nc\nd</p>"); err != nil {
		t.Fatal(err)
	}

	assertSingleEvent(t, w.Body.String(),
		"event: datastar-patch-elements",
		"data: elements <p>a",
		"data: elements b",
		"data: elements c",
		"data: elements d</p>",
	)
}

func assertSingleEvent(t *testing.T, stream string, want ...string) {
	t.Helper()
	if strings.Contains(stream, "\r") {
		t.Errorf("stream contains a carriage return: %q", stream)
	}
	events := clientEvents(stream)
	if len(events) != 1 {
		t.Fatalf("client parsed %d events, want 1: %q", len(events), events)
	}
	if strings.Join(events[0], "\n") != strings.Join(want, "\n") {
		t.Errorf("client parsed %q, want %q", events[0], want)
	}
}

func TestPatchElementsCannotInjectAnEvent(t *testing.T) {
	sse, w := newRecordedSSE(t)
	payload := "hi\r\revent: datastar-patch-signals\rdata: signals {injected: true}\r\r"

	if err := sse.PatchElements(`<p id="post">`+payload+`</p>`,
		WithMode(ElementPatchModeAppend), WithSelectorID("posts")); err != nil {
		t.Fatal(err)
	}

	events := clientEvents(w.Body.String())
	if len(events) != 1 {
		t.Fatalf("client parsed %d events, want 1: %q", len(events), events)
	}
	if events[0][0] != "event: datastar-patch-elements" {
		t.Errorf("first line = %q, want the patch-elements event type", events[0][0])
	}
	for _, line := range events[0][1:] {
		if !strings.HasPrefix(line, "data: ") {
			t.Errorf("client sees a non-data field %q", line)
		}
	}
}

func TestPatchSignalsSplitsOnEverySSELineTerminator(t *testing.T) {
	sse, w := newRecordedSSE(t)

	if err := sse.PatchSignals([]byte("{\"a\":1,\r\"b\":2,\r\n\"c\":3}")); err != nil {
		t.Fatal(err)
	}

	assertSingleEvent(t, w.Body.String(),
		"event: datastar-patch-signals",
		`data: signals {"a":1,`,
		`data: signals "b":2,`,
		`data: signals "c":3}`,
	)
}

func TestSendRejectsLineTerminatorsInSingleLineFields(t *testing.T) {
	for _, terminator := range []string{"\r", "\n", "\r\n"} {
		cases := map[string]func(*ServerSentEventGenerator) error{
			"data line": func(sse *ServerSentEventGenerator) error {
				return sse.Send(EventTypePatchElements, []string{"elements a" + terminator + "b"})
			},
			"event id": func(sse *ServerSentEventGenerator) error {
				return sse.Send(EventTypePatchElements, []string{"elements a"}, WithSSEEventId("1"+terminator+"2"))
			},
			"event type": func(sse *ServerSentEventGenerator) error {
				return sse.Send(EventType("x"+terminator+"y"), []string{"elements a"})
			},
			"selector": func(sse *ServerSentEventGenerator) error {
				return sse.PatchElements("<p></p>", WithSelector("#a"+terminator+"b"))
			},
		}
		for name, send := range cases {
			t.Run(name+" "+strings.NewReplacer("\r", `\r`, "\n", `\n`).Replace(terminator), func(t *testing.T) {
				sse, w := newRecordedSSE(t)
				if err := send(sse); err == nil {
					t.Error("expected an error, got nil")
				}
				if w.Body.Len() != 0 {
					t.Errorf("wrote %q to the stream, want nothing", w.Body.String())
				}
			})
		}
	}
}
