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
export AKASH_HOME="$(realpath ~/.akash)"

export TF_LOG_PROVIDER=DEBUG
```

## Clean terraform
```shell
make clean
```

## Build the provider

Run the following command to build the provider

```shell
# Run the following command to build the provider
go build -o terraform-provider-akash

# Build and install the provider
make install
```

## Development Script

You can run all the commands below by executing:

```shell
make develop
```

## Akash Testing

### Close the Deployment

```shell
./bin/akash tx deployment close --dseq <dseq> --owner $AKASH_ACCOUNT_ADDRESS --from $AKASH_KEY_NAME -y --fees 800uakt --gas auto
```

### Get deployment details

```shell
./bin/akash provider lease-status --home ~/.akash --dseq <dseq> --provider <provider>
```

## Troubleshooting

### `Error: error unmarshalling: invalid character '<' looking for beginning of value`
If you encounter this error close the deployment and try again.
If in development mode, try to increase the fees on deployment creation.