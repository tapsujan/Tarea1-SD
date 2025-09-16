package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Crear router Gin
	r := gin.Default()

	// Endpoint de prueba
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Levantar servidor en puerto 8080
	r.Run(":8080")
}
