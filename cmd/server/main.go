// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// ⚠️ Cargar variables solo si el archivo .env existe (útil para entorno local)
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("⚠️ No se pudo cargar .env: %v", err)
		} else {
			log.Println("✅ Archivo .env cargado exitosamente.")
		}
	} else {
		log.Println("ℹ️ Archivo .env no encontrado, se asume entorno de producción.")
	}

	db.InitCockroach()

	hasData, err := repository.HasStocks()
	if err != nil {
		log.Fatal("❌ Error checking stock data:", err)
	}

	if !hasData {
		log.Println("📥 No stocks found in DB. Importing...")
		if err := service.FetchAndStoreAllStocks(); err != nil {
			log.Fatal("❌ Error fetching stock data:", err)
		}
	} else {
		log.Println("✅ Stock data already present in DB. Skipping import.")
	}

	// Servidor API
	r := mux.NewRouter()
	r.HandleFunc("/stocks", api.GetStocks).Methods("GET")
	r.HandleFunc("/recommendation", api.RecommendMultipleStocks)

	log.Println("✅ Server running on :8080")

	// Obtener origen permitido desde variable de entorno
	origins := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := []string{"http://localhost:5173"} // valor por defecto

	if origins != "" {
		allowedOrigins = splitAndTrim(origins)
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var trimmed []string
	for _, p := range parts {
		trimmed = append(trimmed, strings.TrimSpace(p))
	}
	return trimmed
}
