# akash_deployment Resource

The deployment resource allows you to create a deployment on the Akash Network.

## Example Usage

```terraform
resource "akash_deployment" "my_deployment" {
  sdl = file("sdl.yaml")
  provider_filters {
    providers = ["akash..."]
    enforce = false
  }
}
```

## Argument Reference

The following arguments are required:

- `sdl` - (Required) The SDL configuration of the deployment.

The following arguments are optional:

### provider_filters

- `providers` - (Optional) The list of provider addresses we want to deploy on.
- `enforce` - (Optional) Whether to enforce the filters or to ignore them in case no bid/provider matches the filters.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

### Deployment

- `services` - The services created by the deployment and its URLs.