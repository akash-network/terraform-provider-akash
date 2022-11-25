terraform {
  required_providers {
    akash = {
      version = "0.0.7"
      source  = "Joaos-MacBook-Pro.local/cloud/akash"
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

data "akash_providers" "active" {
  all_providers = false
  minimum_uptime = 100
}

output "all_provider" {
  value = data.akash_providers.active.providers
}