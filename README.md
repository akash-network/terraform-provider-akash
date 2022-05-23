# Terraform Provider Akash

## Development environment setup

```shell
export AKASH_KEY_NAME=terraform
export AKASH_KEYRING_BACKEND=os
export AKASH_ACCOUNT_ADDRESS="$(./bin/akash keys show $AKASH_KEY_NAME -a)"
export AKASH_NET="https://raw.githubusercontent.com/ovrclk/net/master/mainnet"
export AKASH_VERSION="$(curl -s "$AKASH_NET/version.txt")"
export AKASH_CHAIN_ID="$(curl -s "$AKASH_NET/chain-id.txt")"
export AKASH_NODE="http://akash.c29r3.xyz:80/rpc"

export TF_LOG_PROVIDER=DEBUG
```

## Clean terraform
```shell
rm -rf examples/.terraform examples/.terraform.lock.hcl examples/terraform.tfstate examples/terraform.tfstate.backup
```

## Build the provider

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

## Akash Testing



```shell
./bin/akash tx deployment close --dseq <dseq> --owner <owner> --from $AKASH_KEY_NAME -y --fees 5000uakt
```