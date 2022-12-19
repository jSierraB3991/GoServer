package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/OkabeRitarou/GoServer/database"
	"github.com/OkabeRitarou/GoServer/repository"
	"github.com/OkabeRitarou/GoServer/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	JwtSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {

	if config.Port == "" {
		return nil, errors.New("The port is required")
	}

	if config.JwtSecret == "" {
		return nil, errors.New("The secret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("The Database conection is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	handler := cors.Default().Handler(b.router)
	binder(b, b.router)
	repo, err := database.NewPostgreRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	go b.hub.Run()
	repository.SetUserRepository(repo)
	repository.SetPostRepository(repo)

	log.Println("Starting server on port ", b.config.Port)
	err = http.ListenAndServe(b.config.Port, handler)
	if err != nil {
		log.Fatal("ListeneAndServe: ", err)
	}
}
