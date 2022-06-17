terraform {
  required_providers {
    akash = {
      version = "0.3"
      source  = "joaoluna.com/edu/akash"
    }
  }
}

/*
data "akash_deployments" "all" {}

output "all_deployments" {
  value = data.akash_deployments.all.deployments
}

output "all_active_deployments" {
  value = {
    for deployment in data.akash_deployments.all.deployments :
    deployment.deployment_dseq => deployment
    if deployment.deployment_state == "active"
  }
}*/
