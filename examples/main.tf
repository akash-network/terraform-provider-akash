/**
 * Use this file as the starting point for testing this provider.
 **/

terraform {
  required_providers {
    akash = {
      version = "0.0.4"
      source  = "joaoluna.com/cloud/akash"
    }
  }
}

provider "akash" {
}

resource "akash_deployment" "my_deployment" {
  sdl = file("./wordpress.yaml")
  provider_filters {
    providers = ["akash123"]
    enforce = false
  }
}

output "deployment_id" {
  value = akash_deployment.my_deployment.id
}
