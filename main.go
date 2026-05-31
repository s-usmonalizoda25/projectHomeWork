package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/handlers"
	"project/logger"
	"project/middleware"
	"project/storage"

	"go.uber.org/zap"
)

func main() {
	log := logger.New(true)
	defer log.Sync()
	st := &storage.UserStorage{
		FileName: "data/users.json",
		Log:      log,
	}
	h := &handlers.UserHandler{
		Storage: st,
		Log:      log,
	}
	mux := http.NewServeMux()

	mux.Handle("/users", middleware.Logging(log, middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateUser(w, r)
		} else {
			h.GetUsers(w, r)
		}
	}))))

	mux.Handle("/users/{id}", middleware.Logging(log, middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdateUser(w, r)
		} else {
			h.GetUserByID(w, r)
		}
	}))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Сервер успешно запущен", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Критическая ошибка при работе сервера", zap.Error(err))
		}
	}()

	<-shutdownChan
	log.Warn("Получен сигнал остановки!Завершение работы")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Ошибка при остановке сервера", zap.Error(err))
	}

	log.Info("Сервер безопасно остановлен.Данные сохранены!")
}


// Сощдание нового ползователя post
// curl -X POST http://localhost:8080/users \
//   -H "Authorization: secret" \
//   -H "Content-Type: application/json" \
//   -d '{"name": "Suhrob"}'


// Получение списка get
// curl -X GET http://localhost:8080/users \
//   -H "Authorization: secret"


// Ошибка авторизации
// curl -X GET http://localhost:8080/users
