package loghandler

import (
	"context"
	"log/slog"
	"unsafe"

	gaelog "google.golang.org/appengine/v2/log"
)

// LogHandler ...
type LogHandler struct {
	group string
	attrs []slog.Attr
}

// Enabled ...
func (h *LogHandler) Enabled(ctx context.Context, lev slog.Level) bool {
	return true
}

// Handle ...
func (h *LogHandler) Handle(ctx context.Context, rec slog.Record) error {
	buf := make([]byte, 0, 100)
	buf = append(buf, "%s"...)
	args := make([]interface{}, 1, 1+2*(len(h.attrs)+rec.NumAttrs()))
	args[0] = rec.Message
	for _, a := range h.attrs {
		buf = append(buf, " %s=%v"...)
		args = append(args, h.group+a.Key, a.Value)
	}
	rec.Attrs(func(a slog.Attr) bool {
		buf = append(buf, " %s=%v"...)
		args = append(args, h.group+a.Key, a.Value)
		return true
	})
	fmt := unsafe.String(unsafe.SliceData(buf), len(buf))
	switch rec.Level {
	case slog.LevelDebug:
		gaelog.Debugf(ctx, fmt, args...)
	case slog.LevelInfo:
		gaelog.Infof(ctx, fmt, args...)
	case slog.LevelWarn:
		gaelog.Warningf(ctx, fmt, args...)
	case slog.LevelError:
		gaelog.Errorf(ctx, fmt, args...)
	}
	return nil
}

// WithAttrs ...
func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		group: h.group,
		attrs: append(h.attrs[:len(h.attrs):len(h.attrs)], attrs...),
	}
}

// WithGroup ...
func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		group: h.group + name + ".",
		attrs: h.attrs,
	}
}

var _ slog.Handler = (*LogHandler)(nil)
