provider "akash" {
  account_address = "akash..."
  keyring_backend = "os"
  key_name = "terraform"
  node = "http://akash.c29r3.xyz:80/rpc"
  chain_id = "akashnet-2"
  chain_version = "0.16.4"
}