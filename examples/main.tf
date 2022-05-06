terraform {
  required_providers {
    akash = {
      version = "0.3"
      source  = "joaoluna.com/edu/akash"
    }
  }
}

provider "akash" {}

module "akash_deployments" {
  source = "./deployment"
}

output "akash_deployments" {
  value = module.akash_deployments.all_active_deployments
}
