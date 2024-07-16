package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"main.go/cmd/internal/config"
	"main.go/cmd/internal/handlers"
	"main.go/cmd/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.MustLoad()

	// Инициализация логгера
	log := setupLogger(cfg.Env)
	log.Info("starting time_tracker service", slog.String("env", cfg.Env))
	log.Debug("debug message")

	// Инициализация хранилища для аккаунтов
	storage.InitStorage()

	// Регистрация HTTP-обработчиков
	http.HandleFunc("/accounts", handlers.CreateAccountHandler)
	http.HandleFunc("/accounts/", accountOperationsHandler)

	// Настройка HTTP-сервера
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("HTTP сервер запущен на", slog.String("адрес", cfg.HTTPServer.Address))

	// Запуск HTTP-сервера
	err := server.ListenAndServe()
	if err != nil {
		log.Error("Ошибка запуска сервера", slog.String("ошибка", err.Error()))
		os.Exit(1)
	}
}

// setupLogger настраивает логгер в зависимости от окружения
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// accountOperationsHandler обрабатывает операции с аккаунтом
func accountOperationsHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем путь после "/accounts/"
	path := strings.TrimPrefix(r.URL.Path, "/accounts/")

	// Проверяем, какой тип операции требуется
	if strings.HasSuffix(path, "/deposit") {
		if r.Method == http.MethodPost {
			handlers.DepositHandler(w, r)
			return
		}
	} else if strings.HasSuffix(path, "/withdraw") {
		if r.Method == http.MethodPost {
			handlers.WithdrawHandler(w, r)
			return
		}
	} else if strings.HasSuffix(path, "/balance") {
		if r.Method == http.MethodGet {
			handlers.GetBalanceHandler(w, r)
			return
		}
	}

	// Если путь или метод не соответствует ни одной из операций, возвращаем ошибку 404
	http.Error(w, "Not Found", http.StatusNotFound)
}
