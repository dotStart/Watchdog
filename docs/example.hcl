refresh-rate = "5m"
image-dir = "images/"

site "status.example.org" {
  display-name = "Test Site"

  service "web" {
    display-name = "Web Store"
    consul-addr = "127.0.0.1:8700"
    consul-services = [
      "nginx",
      "tomcat"
    ]
  }

  service "db" {
    display-name = "Customer Database"
    consul-addr = "10.240.12.1:8200"
    consul-services = [
      "postgres"
    ]
  }
}

site "status.example.net" {
  display-name = "Test Site"

  service "IRC" {
    display-name = "IRC Network"
    consul-services = [
      "irc"
    ]
  }
}
