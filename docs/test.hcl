refresh-rate = "2m"
image-dir = "docs/images/"

site "localhost" {
  display-name = "Example Site"

  service "a" {
    display-name = "Service A"
    consul-services = [
      "service-a"
    ]
  }
  service "b" {
    display-name = "Service B"
    consul-services = [
      "service-b"
    ]
  }
  service "c" {
    display-name = "Service C"
    consul-services = [
      "service-a",
      "service-b",
    ]
  }
  service "d" {
    display-name = "Service D"
    consul-services = [
      "service-d"
    ]
  }
}
