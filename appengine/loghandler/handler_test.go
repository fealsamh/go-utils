package loghandler

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type secret string

func (secret) LogValue() slog.Value { return slog.StringValue("***") }

func captureFallback(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	old := fallbackOutput
	fallbackOutput = &buf
	t.Cleanup(func() { fallbackOutput = old })
	return &buf
}

func TestHandlerFallback(t *testing.T) {
	require := require.New(t)
	out := captureFallback(t)

	// a non-App Engine context must not panic
	logger := slog.New(new(LogHandler))
	logger.Info("abcd efgh", "a1", "abcd", "a2", 1234)
	require.True(strings.HasSuffix(out.String(), " INFO: abcd efgh a1=abcd a2=1234\n"), out.String())

	tm := time.Date(2026, time.June, 27, 21, 20, 30, 0, time.UTC)
	r := slog.NewRecord(tm, slog.LevelWarn+2, "msg", 0)
	out.Reset()
	require.Nil(new(LogHandler).Handle(context.Background(), r))
	require.Equal("2026-06-27T21:20:30Z WARN+2: msg\n", out.String())
}

func TestHandlerAttrs(t *testing.T) {
	require := require.New(t)
	out := captureFallback(t)

	logger := slog.New(new(LogHandler)).
		With("svc", "parser").
		WithGroup("").
		WithGroup("g1").
		With("n", 1).
		WithGroup("g2")
	logger.Info("msg",
		"s", "two words",
		"e", "",
		"q", `a="b"`,
		"err", errors.New("boom"),
		"pw", secret("hunter2"),
		slog.Group("req", "id", 7, slog.Group("", "inl", true)),
		slog.Attr{},
	)
	line := out.String()
	_, line, _ = strings.Cut(line, " INFO: ")
	require.Equal(`msg svc=parser g1.n=1 g1.g2.s="two words" g1.g2.e="" g1.g2.q="a=\"b\"" g1.g2.err=boom g1.g2.pw=*** g1.g2.req.id=7 g1.g2.req.inl=true`+"\n", line)
}

func TestLevels(t *testing.T) {
	require := require.New(t)
	out := captureFallback(t)

	logger := slog.New(new(LogHandler))
	for _, lev := range []slog.Level{slog.LevelDebug - 4, slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError, LevelCritical} {
		logger.Log(context.Background(), lev, "x")
	}
	require.Equal(6, strings.Count(out.String(), "\n"))
}
