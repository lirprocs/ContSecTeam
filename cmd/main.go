package main

import (
	"ContSecTeam/config"
	"ContSecTeam/internal/handler"
	"ContSecTeam/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.New()
	srv := service.NewService(cfg.QueueSize)
	ctx := context.Background()

	srv.Start(ctx, cfg.Workers)

	h := handler.NewHandler(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/enqueue", h.Enqueue)
	mux.HandleFunc("/healthz", h.Healthz)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("server started on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	log.Println("Получен сигнал завершения, остановка сервера...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	srv.Stop()

	log.Println("Сервер и воркеры успешно остановлены.")
}
