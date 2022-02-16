package utils

import (
	"context"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
)

func GetUserEmail(c context.Context) string {
	return config.GetConfig().Mail.EmaiFrom
}
