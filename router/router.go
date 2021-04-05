package router

import (
	"api/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":8080"

func CreateRouter() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/read", handlers.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/read/{pokemonId}", handlers.GetById).Methods(http.MethodGet)
	api.HandleFunc("/getBerries", handlers.GetBerries).Methods(http.MethodGet)
	api.HandleFunc("/readConcurrently", handlers.ReadConcurrently).Methods(http.MethodGet)

	log.Println("Server started listening on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
