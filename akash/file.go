package akash

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateTemporaryFile creates a temporary file on the filesystem.
func CreateTemporaryFile(ctx context.Context, sdl string) (string, error) {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s/deployment-%d.yaml", os.TempDir(), timestamp)
	tflog.Debug(ctx, fmt.Sprintf("Creating temporary deployment file %s", filename))

	err := ioutil.WriteFile(filename, []byte(sdl), 0666)
	if err != nil {
		return "", err
	}

	return filename, nil
}
