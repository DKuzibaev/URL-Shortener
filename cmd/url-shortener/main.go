package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	logcsv "url-shortener/internal/logs"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv go get -u github.com/ilyakaznacheev/cleanenv
	cfg := config.MustLoad()

	// TODO: init loger: slog
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage: sglline (clickhouse)

	// TODO: init router: chi, "chi render"

	// TODO: run server"
}

func setupLogger(env string) *slog.Logger {
	var handlers []slog.Handler

	// Консольный handler
	switch env {
	case envLocal:
		handlers = append(handlers, slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev, envProd:
		handlers = append(handlers, slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	// CSV handler
	csvHandler := logcsv.NewCSVHandler("logs/log.csv", slog.LevelDebug)
	handlers = append(handlers, csvHandler)

	// MultiHandler
	multi := logcsv.NewMultiHandler(handlers...)

	return slog.New(multi)
}
