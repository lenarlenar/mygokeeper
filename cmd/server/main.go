package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lenarlenar/mygokeeper/internal/server"
	"github.com/lenarlenar/mygokeeper/internal/server/handler"
	"github.com/lenarlenar/mygokeeper/internal/server/migrate"
	"github.com/lenarlenar/mygokeeper/internal/server/repo"
	"github.com/lenarlenar/mygokeeper/internal/server/service"
	_ "github.com/lib/pq"
)

func main() {

	cfg := server.Load()

	migrate.ApplyMigrations("../../internal/server/migrations", cfg.DBConn)
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

	log.Println("Server started at", cfg.ServerAddr)
	err = http.ListenAndServe(cfg.ServerAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
