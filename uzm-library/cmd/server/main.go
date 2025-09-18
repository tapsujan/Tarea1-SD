package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

  "github.com/gin-gonic/gin"
	"github.com/tapsujan/Tarea1-SD/internal/models"
)


func loadSQL() {
	sqlFile, err := ioutil.ReadFile("./migrations/init.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	
	db, err := models.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	_, err = db.DB().Exec(string(sqlFile))
	if err != nil {
		log.Fatalf("Error executing SQL commands: %v", err)
	}

	fmt.Println("Database initialized successfully!")
}

func createBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	db, err := models.ConnectDB()
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database"})
			return
	}
	defer db.Close()

	if result := db.Create(&book); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create book", "details": result.Error.Error()})
			return
	}

	c.JSON(http.StatusCreated, book)
}

func getBooks(c *gin.Context) {
	db, err := models.ConnectDB()
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database"})
			return
	}
	defer db.Close()

	var books []models.Book
	if result := db.Find(&books); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve books"})
			return
	}

	c.JSON(http.StatusOK, books)
}

func main() {
	loadSQL()

	// Crear router Gin
	r := gin.Default()

	// Endpoint de prueba
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/books", createBook)
	r.GET("/books", getBooks)

	// Levantar servidor en puerto 8080
	r.Run(":8080")
}
