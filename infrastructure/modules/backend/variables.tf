variable "name" {
  description = "Base name for the container"
  type        = string
}

variable "image" {
  description = "Docker image to use for the backend"
  type        = string
}

variable "internal_port" {
  description = "Internal port the backend listens on"
  type        = number
}

variable "external_port" {
  description = "Base external port for published containers"
  type        = number
}

variable "replicas" {
  description = "Number of backend replicas to launch"
  type        = number
  default     = 1
}

variable "network_name" {
  description = "Docker network name to connect the container"
  type        = string
}

variable "env_vars" {
  description = "Environment variables for backend container"
  type        = map(string)
  default     = {}
}