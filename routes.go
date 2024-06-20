package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Loading() {

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/login", Login).Methods("POST")

	r.HandleFunc("/signup", Singup).Methods("POST")
	r.HandleFunc("/admin", AdminView).Methods("GET")
	r.HandleFunc("/upload", UploadImage).Methods("POST")
	r.HandleFunc("/delete", DeleteImage).Methods("DELETE")
	r.HandleFunc("/show", ShowAll).Methods("GET")
	r.HandleFunc("/teste", ImageGen).Methods("GET")
	r.HandleFunc("/teste2", teste).Methods("GET")

	// Adicione suporte a CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))

}
