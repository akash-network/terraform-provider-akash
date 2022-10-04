package akash

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/ioutil"
	"os"
	"time"
)

func CreateTemporaryDeploymentFile(ctx context.Context, sdl string) (string, error) {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s/deployment-%d.yaml", os.TempDir(), timestamp)
	tflog.Debug(ctx, fmt.Sprintf("Creating temporary deployment file %s", filename))

	err := ioutil.WriteFile(filename, []byte(sdl), 0666)
	if err != nil {
		return "", err
	}

	return filename, nil
}
