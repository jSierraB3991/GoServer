package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/OkabeRitarou/GoServer/handlers"
	"github.com/OkabeRitarou/GoServer/middleware"
	"github.com/OkabeRitarou/GoServer/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func BindRoutes(s server.Server, r *mux.Router) {
	r.Use(middleware.CheckoutMiddleware(s))
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MyHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts", handlers.ListPostsHandlers(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)

	r.HandleFunc("/ws", s.Hub().HandleWebSocket)
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading file .env")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JwtSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}
	s.Start(BindRoutes)
}
