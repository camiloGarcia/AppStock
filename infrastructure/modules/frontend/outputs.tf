output "container_names" {
  value = [for c in docker_container.frontend : c.name]
}

output "published_ports" {
  value = [for c in docker_container.frontend : c.ports[0].external]
}
