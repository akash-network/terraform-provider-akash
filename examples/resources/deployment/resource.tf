resource "akash_deployment" "my_deployment" {
  sdl = file("./path/to/file.yaml")
  provider_filters {
    providers = ["provider0address1"]
    enforce = false
  }
}