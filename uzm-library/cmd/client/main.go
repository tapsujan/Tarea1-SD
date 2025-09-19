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
	"time"
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
			shoppingCart(userID, reader)
		case "3":
			// Loans() *
		case "4":
			myAccount(userID, reader)
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

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Println("Error al parsear respuesta:", err)
		return
	}

	fmt.Println("\nCatálogo de libros:")
	booksData, ok := respData["books"].([]interface{})
	if !ok || len(booksData) == 0 {
		fmt.Println("No hay libros registrados.")
		return
	}

	// Encabezado
	fmt.Printf("\n%-8s | %-20s | %-10s | %-10s | %-6s\n", "ID libro", "Nombre", "Categoria", "Modalidad", "Valor")
	fmt.Println(strings.Repeat("-", 80))

	for _, b := range booksData {
		book := b.(map[string]interface{})
		fmt.Printf("%-8v | %-20v | %-10v | %-8v | %-6v\n",
			book["ID"],
			book["BookName"],
			book["BookCategory"],
			book["TransactionType"],
			book["Price"],
		)
	}
}

// carrito (falta optimizar)
func shoppingCart(userID uint, reader *bufio.Reader) {
	var books []map[string]uint

	fmt.Println("Ingrese los IDs de los libros a agregar al carro (Enter vacío para terminar):")
	for {
		fmt.Print("Ingrese el ID del libro a agregar al carro: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			break
		}

		var id uint
		_, err := fmt.Sscan(input, &id)
		if err != nil {
			fmt.Println("ID inválido, intente nuevamente")
			continue
		}
		books = append(books, map[string]uint{"id": id})
	}

	if len(books) == 0 {
		fmt.Println("Carro vacío.")
		return
	}

	request := map[string]interface{}{
		"user_id": userID,
		"books":   books,
	}

	data, _ := json.Marshal(request)
	resp, err := http.Post(baseURL+"/cart/checkout", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error al enviar carrito:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", string(body))
		return
	}

	var respData map[string]interface{}
	if err := json.Unmarshal(body, &respData); err != nil {
		fmt.Println("Error al parsear respuesta:", err)
		return
	}

	// Mostrar resumen del carro
	if totalLibros, ok1 := respData["total_books"]; ok1 {
		if totalPrecio, ok2 := respData["total_price"]; ok2 {
			fmt.Printf("\nSu carro es de %v libros para un total de %v usm pesos.\n", totalLibros, totalPrecio)
		}
	}

	if booksData, ok := respData["books"].([]interface{}); ok && len(booksData) > 0 {
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("| %-20s | %-10s | %-6s | %-15s |\n", "Nombre", "Modalidad", "Valor", "Fecha devolución")
		fmt.Println("------------------------------------------------------------")

		for _, b := range booksData {
			book := b.(map[string]interface{})
			name := book["BookName"]
			tipo := book["TransactionType"]
			precio := book["Price"]

			fechaDevolucion := "-"
			if tipo == "Arriendo" {
				// fecha: un mes desde hoy
				fechaDevolucion = time.Now().AddDate(0, 1, 0).Format("02/01/2006")
			}

			fmt.Printf("| %-20v | %-10v | %-6v | %-15v |\n", name, tipo, precio, fechaDevolucion)
		}
		fmt.Println("------------------------------------------------------------")
	}

	fmt.Print("Confirmar pedido: [Enter] ")
	reader.ReadString('\n')

	if msg, ok := respData["message"]; ok {
		fmt.Println(msg)
	}
}

func myAccount(userID uint, reader *bufio.Reader) {
	for {
		fmt.Println("\nMi cuenta")
		fmt.Println("1. Consultar saldo")
		fmt.Println("2. Abonar usm pesos")
		fmt.Println("3. Ver historial de compras y arriendos")
		fmt.Println("4. Salir")
		fmt.Print("Seleccione una opción: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			consultarSaldo(userID)
		case "2":
			abonarSaldo(userID, reader)
		case "3":
			historial(userID)
		case "4":
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func consultarSaldo(userID uint) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, userID))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var user map[string]interface{}
	json.Unmarshal(body, &user)

	fmt.Printf("Su saldo es de %v usm pesos\n", user["UsmPesos"])
}

func abonarSaldo(userID uint, reader *bufio.Reader) {
	fmt.Print("Ingrese la cantidad de usm pesos a abonar: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var monto int
	_, err := fmt.Sscan(input, &monto)
	if err != nil {
		fmt.Println("Monto inválido")
		return
	}

	payload := map[string]interface{}{"abonar": monto}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/users/%d", baseURL, userID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", string(body))
		return
	}
	// nuevo saldo
	var respData map[string]interface{}
	json.Unmarshal(body, &respData)

	if saldo, ok := respData["saldo"]; ok {
		fmt.Printf("Nuevo saldo de %v usm pesos\n", saldo)
	}

}

func historial(userID uint) {
	resp, err := http.Get(fmt.Sprintf("%s/transactions?user_id=%d", baseURL, userID))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var trans []map[string]interface{}
	json.Unmarshal(body, &trans)

	if len(trans) == 0 {
		fmt.Println("No tiene transacciones registradas.")
		return
	}

	respBooks, _ := http.Get(baseURL + "/books")
	var booksResp map[string][]map[string]interface{}
	json.NewDecoder(respBooks.Body).Decode(&booksResp)
	books := booksResp["books"]

	bookMap := make(map[interface{}]string)
	for _, b := range books {
		bookMap[b["ID"]] = b["BookName"].(string)
	}

	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Printf("| %-12s | %-8s | %-20s | %-10s | %-15s | %-6s |\n", "ID transacción", "ID libro", "Nombre", "Tipo", "Fecha", "Valor")
	fmt.Println("-------------------------------------------------------------------------------------------")

	for _, t := range trans {
		bookName := "-"
		if t["BookID"] != nil {
			bookName = bookMap[t["BookID"]]
		}
		fmt.Printf("| %-12v | %-8v | %-20v | %-10v | %-15v | %-6v |\n",
			t["ID"], t["BookID"], bookName, t["Type"], t["Date"], t["Amount"])
	}
	fmt.Println("-------------------------------------------------------------------------------------------")
}
