package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

const (
	timeFormat = "[15:04:05.000]"

	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

type DevHandler struct {
	h slog.Handler
	r func([]string, slog.Attr) slog.Attr
	b *bytes.Buffer
	m *sync.Mutex
}

func NewDevHandler(opts *slog.HandlerOptions) *DevHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}
	}
	if opts.ReplaceAttr == nil {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			return a
		}
	}

	b := &bytes.Buffer{}
	return &DevHandler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		r: opts.ReplaceAttr,
		m: &sync.Mutex{},
	}
}

func (h *DevHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *DevHandler) Handle(ctx context.Context, r slog.Record) error {
	var level string
	levelAttr := slog.Attr{
		Key:   slog.LevelKey,
		Value: slog.AnyValue(r.Level),
	}
	levelAttr = h.r([]string{}, levelAttr)

	if !levelAttr.Equal(slog.Attr{}) {
		level = fmt.Sprintf("%s:", levelAttr.Value.String())

		if r.Level <= slog.LevelDebug {
			level = colorize(lightGray, level)
		} else if r.Level <= slog.LevelInfo {
			level = colorize(green, level)
		} else if r.Level < slog.LevelWarn {
			level = colorize(lightBlue, level)
		} else if r.Level < slog.LevelError {
			level = colorize(lightYellow, level)
		} else if r.Level <= slog.LevelError+1 {
			level = colorize(lightRed, level)
		} else if r.Level > slog.LevelError+1 {
			level = colorize(lightMagenta, level)
		}
	}

	// timestamp
	var timestamp string
	timeAttr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(timeFormat)),
	}
	timeAttr = h.r([]string{}, timeAttr)
	if !timeAttr.Equal(slog.Attr{}) {
		timestamp = colorize(lightGray, timeAttr.Value.String())
	}

	// name
	var name string
	nameAttr := slog.Attr{
		Key:   "name",
		Value: slog.StringValue("gosk"),
	}
	nameAttr = h.r([]string{}, nameAttr)
	if !nameAttr.Equal(slog.Attr{}) {
		name = colorize(blue, nameAttr.Value.String())
	}

	// version
	var version string
	versionAttr := slog.Attr{
		Key:   "version",
		Value: slog.StringValue("0.0.1"),
	}
	versionAttr = h.r([]string{}, versionAttr)
	if !versionAttr.Equal(slog.Attr{}) {
		version = colorize(yellow, versionAttr.Value.String())
	}

	// tags
	var tags string
	tagsAttr := slog.Attr{
		Key:   "tags",
		Value: slog.StringValue("http"),
	}
	tagsAttr = h.r([]string{}, tagsAttr)
	if !tagsAttr.Equal(slog.Attr{}) {
		tags = colorize(lightGray, tagsAttr.Value.String())
	}

	// msg
	var msg string
	msgAttr := slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}
	msgAttr = h.r([]string{}, msgAttr)
	if !msgAttr.Equal(slog.Attr{}) {
		msg = colorize(white, msgAttr.Value.String())
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("log attribute marshaling error: %w", err)
	}

	out := strings.Builder{}
	if len(timestamp) > 0 {
		out.WriteString(timestamp)
		out.WriteString(" ")
	}
	if len(name) > 0 {
		out.WriteString(name)
		out.WriteString(" ")
	}
	if len(version) > 0 {
		out.WriteString(version)
		out.WriteString(" ")
	}
	if len(tags) > 0 {
		out.WriteString(tags)
		out.WriteString(" ")
	}
	if len(level) > 0 {
		out.WriteString(level)
		out.WriteString(" ")
	}
	if len(msg) > 0 {
		out.WriteString(msg)
		out.WriteString(" ")
	}
	if len(bytes) > 2 {
		out.WriteString(colorize(darkGray, string(bytes)))
	}
	fmt.Println(out.String())

	return nil
}

func (h *DevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DevHandler{h: h.h.WithAttrs(attrs), b: h.b, r: h.r, m: h.m}
}

func (h *DevHandler) WithGroup(name string) slog.Handler {
	return &DevHandler{h: h.h.WithGroup(name), b: h.b, r: h.r, m: h.m}
}

func (h *DevHandler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("inner handler error: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("inner handler unmarshal error: %w", err)
	}

	return attrs, nil
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey ||
			a.Key == "name" ||
			a.Key == "version" ||
			a.Key == "tags" ||
			a.Key == "pid" {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}

		return next(groups, a)
	}
}
