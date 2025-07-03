package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lenarlenar/mygokeeper/internal/server"
	"github.com/lenarlenar/mygokeeper/internal/server/handler"
	"github.com/lenarlenar/mygokeeper/internal/server/migrate"
	"github.com/lenarlenar/mygokeeper/internal/server/repo"
	"github.com/lenarlenar/mygokeeper/internal/server/service"
	_ "github.com/lib/pq"
)

func main() {

	cfg := server.Load()

	migrate.ApplyMigrations(cfg.MigrationsPath, cfg.DBConn)
	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := repo.NewPostgresUserRepo(db)
	authService := service.NewAuthService(userRepo)
	authHandler := &handler.AuthHandler{Auth: authService, JWTSecret: cfg.JWTSecret}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	recordRepo := repo.NewPostgresRecordRepo(db)
	recordService := service.NewRecordService(recordRepo)
	recordHandler := &handler.RecordHandler{Service: recordService}

	mux.Handle("/record", handler.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(recordHandler.Save)))
	mux.Handle("/records", handler.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(recordHandler.GetAll)))

	mux.Handle("/record/", handler.AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			recordHandler.Delete(w, r)
		case http.MethodPut:
			recordHandler.Update(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: mux,
	}

	go func() {
		log.Println("Server started at", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
