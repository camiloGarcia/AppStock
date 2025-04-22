resource "docker_network" "appstock_network" {
  name = "appstock_net"
}

locals {
  frontend_url = "http://localhost:3000"
  backend_url  = "http://localhost:9000"
}

module "backend" {
  source         = "./modules/backend"
  name           = "appstock"
  image          = "appstock-backend:1.0.0"
  internal_port  = 8080
  external_port  = 9000
  replicas       = 2
  network_name   = docker_network.appstock_network.name

  env_vars = merge(
    var.backend_env_vars,
    {
      ALLOWED_ORIGINS = local.frontend_url
    }
  )

}

module "frontend" {
  source         = "./modules/frontend"
  name           = "appstock"
  image          = "appstock-frontend:1.0.0"
  internal_port  = 80
  external_port  = 3000
  replicas       = 1
  network_name   = docker_network.appstock_network.name

  env_vars = {
    VITE_API_BASE_URL = local.backend_url
  }
}
