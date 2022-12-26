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

- `services` - The services created by the deployment. This attribute is a list of services with the following structure:
```hcl
services: [
  {
    "available": 1,
    "available_replicas": 1,
    "name": "db",
    "ready_replicas": 1,
    "replicas": 1,
    "total": 1,
    "updated_replicas": 1,
    "uris": ["url1", "url2"],
    "ips": [
      {
        "ip": "xxx.xxx.xxx.xxx",
        "port": 1234,
        "proto": "TCP",
        "external_port": 12345
      }
    ],
    "forwarded_ports": [
      {
        "host": "hosturl",
        "port": 1234,
        "proto": "TCP",
        "external_port": 12345
      }
    ]
  }
]
  ```