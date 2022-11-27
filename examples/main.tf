/**
 * Use this file as the starting point for testing this provider.
 **/

terraform {
  required_providers {
    akash = {
      version = "0.0.5"
      source  = "cloud-j-luna/akash"
    }
  }
}

provider "akash" {
  # Also remember these values can be provided as env variables.
  account_address = "<address>"
  keyring_backend = "os"
  key_name = "terraform"
  node = "http://akash.c29r3.xyz:80/rpc"
  chain_id = "akashnet-2"
  chain_version = "0.16.4"
}

resource "akash_deployment" "my_deployment" {
  sdl = file("./wordpress.yaml")
  provider_filters {
    providers = ["akashpreferredprovider"]
    enforce = false
  }
}

output "deployment_id" {
  value = akash_deployment.my_deployment.id
}
