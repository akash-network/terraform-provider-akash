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

resource "akash_deployment" "my_deployment" {
  sdl = file("./wordpress.yaml")
}

output "my_deployment" {
  value = akash_deployment.my_deployment
}

/*output "akash_deployments" {
  value = module.akash_deployments.all_active_deployments
}*/
