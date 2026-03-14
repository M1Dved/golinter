package example

import "log/slog"

func exampleFunction() {
	slog.Info("запуск сервера")
	slog.Info("starting server")

	slog.Info("Starting server on port 8080")
	slog.Info("starting server on port 8080")

	slog.Info("server started!")
	slog.Info("server started")

	slog.Info("user password is wrong")
	slog.Info("user authenticated")
}
