terraform {
  required_providers {
    akash = {
      version = "0.0.4"
      source  = "cloud-j-luna/akash"
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
  path = "/Users/joaoluna/Documents/Programming/terraform-akash-provider/bin/akash"
}

data "akash_deployments" "all" {}

output "akash_deployments" {
  value = data.akash_deployments.all
}