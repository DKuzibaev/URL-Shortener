package logs

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"log/slog"
)

type CSVHandler struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
	level  slog.Level
}

func NewCSVHandler(filePath string, level slog.Level) *CSVHandler {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("cannot create directory for CSV log: %v", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("cannot open CSV log file: %v", err)
	}

	writer := csv.NewWriter(file)

	// Записать заголовок, если файл пустой
	info, _ := file.Stat()
	if info.Size() == 0 {
		writer.Write([]string{"timestamp", "level", "message", "extra"})
		writer.Flush()
	}

	return &CSVHandler{
		file:   file,
		writer: writer,
		level:  level,
	}
}

func (h *CSVHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CSVHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	extra := ""
	r.Attrs(func(a slog.Attr) bool {
		extra += a.Key + "=" + a.Value.String() + ";"
		return true
	})

	err := h.writer.Write([]string{
		r.Time.Format(time.RFC3339),
		r.Level.String(),
		r.Message,
		extra,
	})
	if err != nil {
		return err
	}

	h.writer.Flush()
	return nil
}

func (h *CSVHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *CSVHandler) WithGroup(name string) slog.Handler       { return h }

func (h *CSVHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.writer.Flush()
	return h.file.Close()
}
