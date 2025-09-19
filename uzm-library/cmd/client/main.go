package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var baseURL = "http://localhost:8080"

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nMenu")
		fmt.Println("1. Registrarse")
		fmt.Println("2. Iniciar sesión")
		fmt.Println("3. Terminar ejecución")
		fmt.Print("Seleccione una opción: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			register(reader)
		case "2":
			login(reader)
		case "3":
			fmt.Println("Muchas gracias por visitarnos")
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func secondMenu(userID uint, reader *bufio.Reader) {
	for {
		fmt.Println("\nMenu")
		fmt.Println("1. Ver catálogo")
		fmt.Println("2. Carro de compras")
		fmt.Println("3. Mis préstamos")
		fmt.Println("4. Mi cuenta")
		fmt.Println("5. Populares")
		fmt.Println("6. Salir")
		fmt.Print("Seleccione una opción: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			catalogo()
		case "2":
			// shoppinCart() *por hacer
		case "3":
			// Loans() *
		case "4":
			// myAccount() *
		case "5":
			// Popular() *
		case "6":
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func register(reader *bufio.Reader) {
	fmt.Print("Ingrese su nombre: ")
	firstName, _ := reader.ReadString('\n')
	fmt.Print("Ingrese su apellido: ")
	lastName, _ := reader.ReadString('\n')
	fmt.Print("Ingrese su email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Ingrese su contraseña: ")
	password, _ := reader.ReadString('\n')

	user := map[string]interface{}{
		"FirstName": strings.TrimSpace(firstName),
		"LastName":  strings.TrimSpace(lastName),
		"Email":     strings.TrimSpace(email),
		"Password":  strings.TrimSpace(password),
		"UsmPesos":  0,
	}

	data, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("\nUsuario creado con exito")
}

func login(reader *bufio.Reader) {
	fmt.Print("Ingrese su email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Ingrese su contraseña: ")
	password, _ := reader.ReadString('\n')

	credentials := map[string]string{
		"email":    strings.TrimSpace(email),
		"password": strings.TrimSpace(password),
	}

	data, _ := json.Marshal(credentials)
	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var respData map[string]interface{}
	json.Unmarshal(body, &respData)

	if resp.StatusCode == 200 {
		user, ok := respData["user"].(map[string]interface{})
		if !ok || user == nil {
			fmt.Println("Credenciales incorrectas o usuario no existe")
			return
		}
		userID := uint(user["ID"].(float64))
		fmt.Println("Sesion iniciada correctamente")
		secondMenu(userID, reader)
	} else {
		fmt.Println("Error:", resp.Status)
	}
}

func catalogo() {
	resp, err := http.Get(baseURL + "/books")
	if err != nil {
		fmt.Println("Error al obtener catálogo:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	var books []map[string]interface{}
	err = json.Unmarshal(body, &books)
	if err != nil {
		fmt.Println("Error al parsear respuesta:", err)
		return
	}

	fmt.Println("\nCatálogo de libros:")
	if len(books) == 0 {
		fmt.Println("No hay libros registrados.")
		return
	}

	// Encabezado
	fmt.Printf("\n%-8s | %-20s | %-10s | %-10s | %-6s\n", "ID libro", "Nombre", "Categoria", "Modalidad", "Valor")
	fmt.Println(strings.Repeat("-", 80))

	for _, b := range books {
		fmt.Printf("%-8v | %-20v | %-10v | %-8v | %-6v\n",
			b["ID"],
			b["BookName"],
			b["BookCategory"],
			b["TransactionType"],
			b["Price"],
		)
	}
}
