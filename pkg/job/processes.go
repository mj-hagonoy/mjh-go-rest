package job

import (
	"context"
	"fmt"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/file_storage"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/user"
)

func ProcessJob(ctx context.Context, job Job) error {
	switch job.Type {
	case JOB_TYPE_IMPORT_USERS:
		return importUsersFromCsv(ctx, job)
	default:
		return fmt.Errorf("unsupported job type %s", job.Type)
	}
}

func importUsersFromCsv(ctx context.Context, job Job) error {
	storage, err := file_storage.GetStorage()
	if err != nil {
		return err
	}
	byteData, err := storage.Read(ctx, job.SourceFile)
	if err != nil {
		return err
	}
	if err := user.ImportUsersFromCsvBytes(ctx, byteData); err != nil {
		return err
	}

	job.SetStatus(JOB_STATUS_DONE)
	return job.Update(ctx)
}
