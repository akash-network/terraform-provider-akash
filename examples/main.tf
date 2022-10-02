/**
 * Use this file as the starting point for testing this provider.
 **/

terraform {
  required_providers {
    akash = {
      version = "0.0.5"
      source  = "joaoluna.com/cloud/akash"
    }
  }
}

provider "akash" {
  account_address = "akash1qyfg4zl2dku8ry7gjkhf88vnc3zrn6vmnzlvr9"
  keyring_backend = "os"
  key_name = "terraform"
  node = "http://akash.c29r3.xyz:80/rpc"
  chain_id = "akashnet-2"
  chain_version = "0.16.4"
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
