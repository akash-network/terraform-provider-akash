# akash Provider

The Akash provider is used to interact with the Akash Network.

## Example Usage

```terraform
provider "akash" {}
```

Remember to set all the required variables for the provider to work properly:

| Variable                | Description                                                      |
|-------------------------|------------------------------------------------------------------|
| `AKASH_KEY_NAME`        | Name of your keychain.                                           |
| `AKASH_KEYRING_BACKEND` | Backend of the keyring.                                          |
| `AKASH_ACCOUNT_ADDRESS` | Address of your account.                                         |
| `AKASH_NET`             | Network to use, usually the mainnet.                             |
| `AKASH_VERSION`         | Version of the network.                                          |
| `AKASH_CHAIN_ID`        | Chain id of the network.                                         |
| `AKASH_NODE`            | Akash node to connect to.                                        |
| `AKASH_HOME`            | Absolute path to the Akash's home folder, usually under ~/.akash |
| `AKASH_PATH`            | (Optional) The path to the Akash binary                          |