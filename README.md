# 📊 AppStock

AppStock es una aplicación fullstack que consume información bursátil desde una API pública, la almacena en una base de datos CockroachDB, y la expone a través de una API REST paginada. Además, incluye una interfaz frontend moderna construida con Vue 3 + TypeScript + Tailwind CSS.

---

## 🧩 Tecnologías utilizadas

### 🔧 Backend (Go)
- Golang 1.21
- API REST con Gorilla Mux
- Base de datos: CockroachDB Cloud
- Control de variables de entorno con `godotenv`
- Paginación y control de CORS

### 🎨 Frontend (Vue 3)
- Vite + Vue 3 + TypeScript
- Tailwind CSS
- Fetch API con paginación dinámica
- Separación en componentes

---

## 🚀 Cómo ejecutar el proyecto

### Backend

```bash
cd AppStock
go run ./cmd/server
