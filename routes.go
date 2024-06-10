package main

import (
	"github.com/gin-gonic/gin"
)

func reAlu(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": "1",
	})
}

func Loading() {

	r := gin.Default()
	r.GET("/", reAlu)

	r.Run()

	// r := mux.NewRouter()
	// r.HandleFunc("/", Home)
	// r.HandleFunc("/login", Login).Methods("POST")

	// r.HandleFunc("/signup", Singup).Methods("POST")
	// r.HandleFunc("/admin", AdminView).Methods("GET")
	// r.HandleFunc("/upload", Upload).Methods("POST")
	// r.HandleFunc("/generate", GenerateImage).Methods("GET")

	// // Adicione suporte a CORS
	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	// origins := handlers.AllowedOrigins([]string{"*"})

	// log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))

}
