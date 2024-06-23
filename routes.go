package main

import (
	"log"
	"net/http"
	"projetov2/Backend-projeto/utility"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Loading() {

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/signup", Singup).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	// r.Handle("/admin", utility.AuthMiddleware(http.HandlerFunc(AdminView))).Methods("GET")
	r.HandleFunc("/admin", AdminView).Methods("GET")
	// r.HandleFunc("/upload", UploadImage).Methods("POST")
	r.Handle("/upload", utility.AuthMiddleware(http.HandlerFunc(UploadImage))).Methods("POST")

	// r.HandleFunc("/delete", DeleteImage).Methods("DELETE")
	r.Handle("/delete", utility.AuthMiddleware(http.HandlerFunc(DeleteImage))).Methods("DELETE")

	// r.HandleFunc("/show", ShowAll).Methods("GET")
	r.HandleFunc("/teste", ImageGen).Methods("GET")
	// r.Handle("/teste", utility.AuthMiddleware(http.HandlerFunc(ImageGen))).Methods("GET")

	// r.HandleFunc("/teste2", teste).Methods("GET")
	r.Handle("/teste2", utility.AuthMiddleware(http.HandlerFunc(teste))).Methods("GET")

	// Adicione suporte a CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	credentials := handlers.AllowCredentials()

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins, credentials)(r)))

}
