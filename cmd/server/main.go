// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"appstock/internal/api"
	"appstock/internal/repository"
	"appstock/internal/service"
	"appstock/pkg/db"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Obtener el directorio del ejecutable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error obteniendo la ruta del ejecutable: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Construir la ruta completa al archivo .env
	envPath := filepath.Join(exeDir, ".env")

	// Cargar las variables de entorno desde el archivo .env
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("❌ Error loading .env file")
	// }

	db.InitCockroach()

	// Verificar si hay stocks ya guardados
	hasData, err := repository.HasStocks()
	if err != nil {
		log.Fatal("❌ Error checking stock data:", err)
	}

	if !hasData {
		log.Println("📥 No stocks found in DB. Importing...")
		err := service.FetchAndStoreAllStocks()
		if err != nil {
			log.Fatal("❌ Error fetching stock data:", err)
		}
	} else {
		log.Println("✅ Stock data already present in DB. Skipping import.")
	}

	// Servidor API
	r := mux.NewRouter()
	r.HandleFunc("/stocks", api.GetStocks).Methods("GET")

	log.Println("✅ Server running on :8080")
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}
