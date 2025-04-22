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

### 1. ⚙️ Configurar variables sensibles

Crea un archivo `backend.tfvars` en la carpeta `infrastructure/` con contenido similar a:

```hcl
ALLOWED_ORIGINS = "http://localhost:3000"
STOCK_API_URL   = "{dominio}/production/swechallenge/list"
STOCK_API_KEY   = "TU_API_KEY"
CONN_STR        = "postgresql://usuario:contraseña@servidor.cockroachlabs.cloud:port/DB?sslmode=verify-full"
```


---

### 2. 🔨 Construir imágenes Docker

Desde la raíz del proyecto:

```bash
# Backend
docker build -t appstock-backend:1.0.0 .

# Frontend (pasando la URL del backend como build-arg)
docker build -t appstock-frontend:1.0.0 --build-arg VITE_API_BASE_URL=http://localhost:9000 ./appstock-ui
```

---

### 3. 🧱 Desplegar con Terraform

```bash
cd infrastructure
terraform init
terraform apply -var-file=backend.tfvars
```

---

### 🌐 Acceder a la aplicación

- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend API: [http://localhost:9000/stocks](http://localhost:9000/stocks)

---

## ✅ Estado actual

- ✔️ Consulta de stocks paginada
- ✔️ Filtro por texto
- ✔️ Visualización de recomendaciones por fecha
- ✔️ Despliegue automatizado con Terraform y Docker

---

## 📝 Licencia

MIT