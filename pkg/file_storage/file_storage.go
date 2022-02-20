package file_storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
)

type IFile interface {
	Write(context.Context, string, []byte) error
	Read(context.Context, string) ([]byte, error)
}

var (
	gcp     *GoogleCloudStorage
	gcpOnce sync.Once
)

func GetStorage() (IFile, error) {
	switch config.GetConfig().FileStorage.Default {
	case GOOGLE_CLOUD:
		var gcpErr error
		gcpOnce.Do(func() {
			gcp, gcpErr = NewGoogleCloudClient(
				config.GetConfig().FileStorage.GoogleCloud.ProjectID,
				config.GetConfig().FileStorage.GoogleCloud.BucketName,
				config.GetConfig().FileStorage.GoogleCloud.UploadPath,
			)
		})
		return gcp, gcpErr

	default:
		return nil, fmt.Errorf("GetStorage: unknown type %v", config.GetConfig().FileStorage.Default)
	}
}
