# akash_providers Data Source

The `akash_providers` data source allows you to query the network providers' information.
This data can then be used on the `akash_Deployment`'s provider filters

## Example Usage

```hcl
data "akash_providers" "europe_providers" {
  all_providers = false
  minimum_uptime = 99
  required_attributes = {
    region = "eu-west"
  }
}
```

You can the use the output like this:

```hcl
resource "akash_deployment" "my_deployment" {
  sdl = templatefile("${path.module}/sdl.yaml", {
    // ...
  })
  provider_filters {
    providers = data.akash_providers.europe_providers.providers.*.address // Select all addresses of the providers from the datasource.
  }
}
```

## Argument Reference

The following arguments are available:

- `all_providers` - (Optional) Whether you want to show every provider record available or just the active ones (default is false).
- `minimum_uptime` - (Optional) The minimum uptime the providers must have. Defaults to 0. The higher the value (0-100) the more strict you are on uptime.
- `required_attributes` - (Optional) Key/value pairs that should be present on the provider's attributes. Examples: `arch`, `region`, ...

~> NOTE: There currently isn't a standard in the format of the providers' attributes.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

### `providers`

A `providers` argument is exported which has the following structure:

- `address` - The address of the provider.
- `active` - Whether the provider is active or not.
- `uptime` - The percentage of uptime of the provider.
- `attributes` - Key/value pairs containing the attribute names and values of the provider.
