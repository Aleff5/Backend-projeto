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

	r.GET("/", Home)
	r.POST("/login", Login)
	r.POST("/signup", Singup)
	r.GET("/admin", AdminView)

	r.Run()
	// r := mux.NewRouter()

	// Adicione suporte a CORS
	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	// origins := handlers.AllowedOrigins([]string{"*"})

	// log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))

}
