package example

import "log/slog"

func exampleFunction() {
	slog.Info("запуск сервера") // want "log message should be in English only"
	slog.Info("starting server")

	slog.Info("Starting server on port 8080") // want "log message should start with a lowercase letter"
	slog.Info("starting server on port 8080")

	slog.Info("server started!") // want "log message should not contain special characters or emoji"
	slog.Info("server started")

	slog.Info("user password is wrong") // want "log message may contain sensitive data"
	slog.Info("user authenticated")
}
