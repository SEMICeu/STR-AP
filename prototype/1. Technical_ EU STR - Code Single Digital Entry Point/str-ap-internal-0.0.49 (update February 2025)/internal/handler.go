package internal

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l    *log.Logger
	opts PrettyHandlerOptions
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelError || r.Level >= slog.LevelWarn {
		return h.processRecord(r)
	}

	if r.Level >= h.opts.SlogOpts.Level.Level() {
		return h.processRecord(r)
	}

	return nil
}

func (h *PrettyHandler) processRecord(r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
		opts:    opts,
	}

	return h
}

func getLogLevelFromEnv() slog.Level {
	levelStr := viper.GetString("LOG_LEVEL")
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func CallNewLogger() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: getLogLevelFromEnv(),
		},
	}

	logHandler := NewPrettyHandler(os.Stdout, opts)
	return slog.New(logHandler)
}
