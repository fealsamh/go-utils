package loghandler

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode"

	"github.com/fealsamh/go-utils/nocopy"
	gaelog "google.golang.org/appengine/v2/log"
)

var bs = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// fallbackOutput receives records logged with a context that isn't an App Engine request context.
var fallbackOutput io.Writer = os.Stderr

// LevelCritical is logged as critical in App Engine.
const LevelCritical = slog.LevelError + 4

// LogHandler ...
type LogHandler struct {
	group string // prefix of open groups, e.g. "g1.g2."
	pre   string // attributes added via WithAttrs, already formatted
}

// Enabled ...
func (h *LogHandler) Enabled(ctx context.Context, lev slog.Level) bool {
	return true
}

// Handle ...
func (h *LogHandler) Handle(ctx context.Context, rec slog.Record) error {
	buf := bs.Get().(*bytes.Buffer)
	defer bs.Put(buf)
	buf.Reset()
	buf.WriteString(rec.Message)
	buf.WriteString(h.pre)
	rec.Attrs(func(a slog.Attr) bool {
		appendAttr(buf, h.group, a)
		return true
	})
	if !logToAppEngine(ctx, rec.Level, nocopy.String(buf.Bytes())) {
		line := make([]byte, 0, buf.Len()+40)
		line = rec.Time.AppendFormat(line, time.RFC3339)
		line = append(line, ' ')
		line = append(line, rec.Level.String()...)
		line = append(line, ": "...)
		line = append(line, buf.Bytes()...)
		line = append(line, '\n')
		_, err := fallbackOutput.Write(line)
		return err
	}
	return nil
}

// logToAppEngine logs msg via the App Engine API and reports whether ctx was an App Engine context.
// The API formats msg before returning, so msg may alias a buffer that is reused afterwards.
func logToAppEngine(ctx context.Context, lev slog.Level, msg string) (ok bool) {
	// the App Engine API panics when ctx isn't derived from an App Engine request
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	switch {
	case lev >= LevelCritical:
		gaelog.Criticalf(ctx, "%s", msg)
	case lev >= slog.LevelError:
		gaelog.Errorf(ctx, "%s", msg)
	case lev >= slog.LevelWarn:
		gaelog.Warningf(ctx, "%s", msg)
	case lev >= slog.LevelInfo:
		gaelog.Infof(ctx, "%s", msg)
	default:
		gaelog.Debugf(ctx, "%s", msg)
	}
	return true
}

func appendAttr(buf *bytes.Buffer, prefix string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		if a.Key != "" {
			prefix += a.Key + "."
		}
		for _, ga := range a.Value.Group() {
			appendAttr(buf, prefix, ga)
		}
		return
	}
	buf.WriteByte(' ')
	buf.WriteString(prefix)
	buf.WriteString(a.Key)
	buf.WriteByte('=')
	s := a.Value.String()
	if needsQuoting(s) {
		buf.WriteString(strconv.Quote(s))
	} else {
		buf.WriteString(s)
	}
}

func needsQuoting(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if r == '=' || r == '"' || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

// WithAttrs ...
func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	var buf bytes.Buffer
	buf.WriteString(h.pre)
	for _, a := range attrs {
		appendAttr(&buf, h.group, a)
	}
	return &LogHandler{
		group: h.group,
		pre:   buf.String(),
	}
}

// WithGroup ...
func (h *LogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &LogHandler{
		group: h.group + name + ".",
		pre:   h.pre,
	}
}

var _ slog.Handler = (*LogHandler)(nil)
