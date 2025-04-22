output "container_names" {
  value = [for c in docker_container.backend : c.name]
}

output "published_ports" {
  value = [for c in docker_container.backend : c.ports[0].external]
}
