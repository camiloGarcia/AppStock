terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.2"
    }
  }
}

resource "docker_container" "backend" {
  count = var.replicas

  name  = "${var.name}-backend-${count.index}"
  image = var.image

  ports {
    internal = var.internal_port
    external = var.external_port + count.index
  }

  restart = "always"

  networks_advanced {
    name = var.network_name
  }

  env = [for key, value in var.env_vars : "${key}=${value}"]
}
