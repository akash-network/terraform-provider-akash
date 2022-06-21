package akash

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/ioutil"
)

func CreateTemporaryDeploymentFile(ctx context.Context, sdl string) (string, error) {
	tflog.Debug(ctx, "Creating temporary deployment file /var/tmp/deployment.yaml")

	err := ioutil.WriteFile("/var/tmp/deployment.yaml", []byte(sdl), 0666)
	if err != nil {
		return "", err
	}

	return "/var/tmp/deployment.yaml", nil
}