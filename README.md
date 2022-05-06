# Terraform Provider Hashicups

Run the following command to build the provider

```shell
go build -o terraform-provider-akash
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples && terraform init && terraform apply --auto-approve
```

## Clean terraform
```shell
rm -rf examples/.terraform examples/.terraform.lock.hcl
```

