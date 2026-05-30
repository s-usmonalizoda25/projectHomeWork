package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"project/handlers"
	"project/middleware"
	"project/storage"
)

func main() {
	st := &storage.UserStorage{
		FileName: "data/users.json",
	}
	h := &handlers.UserHandler{
		Storage: st,
	}
	mux := http.NewServeMux()
	mux.Handle("/users", middleware.Logging(middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateUser(w, r)
		} else {
			h.GetUsers(w, r)
		}
	}))))
	mux.Handle("/users/{id}", middleware.Logging(middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Сервер успешно запущен")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка сервера: %v\n", err)
		}
	}()
	<-shutdownChan
	fmt.Println("Получен сигнал остановки!Завершение работы")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при плавной остановке: %v\n", err)
	}
	fmt.Println("Сервер безопасно остановлен.Данные сохранены!")
}