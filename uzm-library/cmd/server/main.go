package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

	c.JSON(http.StatusOK, gin.H{"books": books})
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

	c.JSON(http.StatusOK, gin.H{"users": users})
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

	db, _ := models.ConnectDB()
	defer db.Close()

	var user models.User
	if result := db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// abonar saldo
	if val, ok := input["abonar"]; ok {
		monto := int(val.(float64))
		user.UsmPesos += monto
		db.Save(&user)
		c.JSON(http.StatusOK, gin.H{"message": "Saldo abonado", "saldo": user.UsmPesos})
		return
	}

	if result := db.Model(&user).Updates(input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado", "saldo": user.UsmPesos})
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

type CartRequest struct {
	UserID uint              `json:"user_id"`
	Books  []map[string]uint `json:"books"`
}

// Carrito de compras
func checkoutCart(c *gin.Context) {
	var req CartRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	// Buscar usuario
	var user models.User
	if result := db.First(&user, req.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var selectedBooks []models.Book
	total := 0

	for _, b := range req.Books {
		var book models.Book
		if result := db.Preload("Inventory").First(&book, b["id"]); result.Error != nil {
			continue // libro no encontrado
		}
		if book.Inventory.AvailableQuantity > 0 {
			selectedBooks = append(selectedBooks, book)
			total += book.Price
		}
	}

	if len(selectedBooks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay libros válidos en el carrito"})
		return
	}

	// Caso 1: saldo suficiente
	if user.UsmPesos >= total {
		processCart(db, &user, selectedBooks)
		c.JSON(http.StatusOK, gin.H{
			"message":     "Pedido realizado con éxito",
			"books":       selectedBooks,
			"total_books": len(selectedBooks),
			"total_price": total,
		})
		return
	}

	// Caso 2: optimizar carrito
	if user.UsmPesos > 0 {
		// ordenar libros por precio ascendente
		sort.Slice(selectedBooks, func(i, j int) bool {
			return selectedBooks[i].Price < selectedBooks[j].Price
		})

		var optimized []models.Book
		sum := 0
		for _, book := range selectedBooks {
			if sum+book.Price <= user.UsmPesos {
				optimized = append(optimized, book)
				sum += book.Price
			}
		}

		if len(optimized) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Carrito optimizado por falta de saldo",
				"books":       optimized,
				"total_books": len(optimized),
				"total_price": sum,
				"optimized":   true,
			})
			return
		}

	}

	// Caso 3: no alcanza para nada
	c.JSON(http.StatusBadRequest, gin.H{"error": "Fondos insuficientes"})
}

// procesar carrito: actualiza inventario, saldo y crea registros
func processCart(db *gorm.DB, user *models.User, books []models.Book) {
	for _, book := range books {
		// Descontar inventario
		db.Model(&models.Inventory{}).Where("book_id = ?", book.ID).
			Update("available_quantity", gorm.Expr("available_quantity - ?", 1))

		// Descontar saldo
		user.UsmPesos -= book.Price
		db.Save(user)

		// Popularidad
		db.Model(&book).Update("popularity_score", gorm.Expr("popularity_score + ?", 1))

		// Registrar transacción
		if book.TransactionType == "Venta" {
			sale := models.Sale{UserID: user.ID, BookID: book.ID, SaleDate: time.Now().Format("02/01/2006")}
			db.Create(&sale)

			tx := models.Transaction{UserID: user.ID, BookID: &book.ID, Type: "Venta", Date: sale.SaleDate, Amount: book.Price}
			db.Create(&tx)
		} else if book.TransactionType == "Arriendo" {
			start := time.Now()
			due := start.AddDate(0, 1, 0) // un mes de plazo
			loan := models.Loan{
				UserID:    user.ID,
				BookID:    book.ID,
				StartDate: start.Format("02/01/2006"),
				DueDate:   due.Format("02/01/2006"),
				Status:    "Pendiente",
			}
			db.Create(&loan)

			tx := models.Transaction{UserID: user.ID, BookID: &book.ID, Type: "Arriendo", Date: loan.StartDate, Amount: book.Price}
			db.Create(&tx)
		}
	}
}

// historial
func getTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	db, err := models.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la BD"})
		return
	}
	defer db.Close()

	var transactions []models.Transaction
	if result := db.Where("user_id = ?", userID).Find(&transactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener historial"})
		return
	}

	c.JSON(http.StatusOK, transactions)
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
	r.POST("/cart/checkout", checkoutCart)
	r.GET("/transactions", getTransactions)

	// Levantar servidor en puerto 8080
	r.Run(":8080")
}
