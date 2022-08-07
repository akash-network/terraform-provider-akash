/**
 * Use this file as the starting point for testing this provider.
 **/

terraform {
  required_providers {
    akash = {
      version = "0.0.3"
      source  = "cloud-j-luna/akash"
    }
  }
}

provider "akash" {
}

resource "akash_deployment" "my_deployment" {
  sdl = file("./wordpress.yaml")
  provider_filters {
    provider_preferred = "none"
    enforce = true
  }
}

output "deployment_id" {
  value = akash_deployment.my_deployment.id
}
