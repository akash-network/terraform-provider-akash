# akash_deployment Resource

The deployment resource allows you to create a deployment on the Akash Network.

## Example Usage

```terraform
resource "akash_deployment" "my_deployment" {
  sdl = file("sdl.yaml")
}
```

## Argument Reference

- `sdl` - (Required) The SDL configuration of the deployment.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

### Deployment

- `services` - The services created by the deployment and its URLs.