terraform {
  required_providers {
    akash = {
      version = "0.0.7"
      source  = "cloud-j-luna/akash"
    }
  }
}

provider "akash" {
  account_address = "<address>"
  keyring_backend = "os"
  key_name = "terraform"
  node = "http://akash.c29r3.xyz:80/rpc"
  chain_id = "akashnet-2"
  chain_version = "0.16.4"
}

data "akash_deployments" "all" {}

output "akash_deployments" {
  value = data.akash_deployments.all
}