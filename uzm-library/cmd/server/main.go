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

// Crear usuario
func createUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear usuario"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Listar usuarios
func getUsers(c *gin.Context) {
	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	var users []models.User
	if result := db.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener usuarios"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Obtener usuario por ID
func getUserByID(c *gin.Context) {
	id := c.Param("id")

	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	var user models.User
	if result := db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Actualizar usuario
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	if result := db.Model(&models.User{}).Where("id = ?", id).Updates(input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado"})
}

// Login
func login(c *gin.Context) {
	var input map[string]string

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := input["email"]
	password := input["password"]

	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	var user models.User
	if result := db.Where("email = ? AND password = ?", email, password).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login exitoso", "user": user})
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
	r.POST("/users", createUser)
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUserByID)
	r.PATCH("/users/:id", updateUser)
	r.POST("/login", login)

	// Levantar servidor en puerto 8080
	r.Run(":8080")
}
