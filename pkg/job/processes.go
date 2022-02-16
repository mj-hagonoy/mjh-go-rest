package job

import (
	"context"
	"fmt"

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
	if err := user.ImportUsersFromCsv(ctx, job.SourceFile); err != nil {
		return err
	}

	job.SetStatus(JOB_STATUS_DONE)
	return job.Update(ctx)
}
