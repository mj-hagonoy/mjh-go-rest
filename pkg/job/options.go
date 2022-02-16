package job

import "fmt"

type Option func(*Job) error

func ID(id string) Option {
	return func(b *Job) error {
		b.ID = id
		return nil
	}
}

func Type(t string) Option {
	return func(b *Job) error {
		b.Type = t
		return nil
	}
}

func SourceFile(filepath string) Option {
	return func(b *Job) error {
		b.SourceFile = filepath
		return nil
	}
}

func InitiatedBy(userId string) Option {
	return func(b *Job) error {
		b.InitiatedBy = userId
		return nil
	}
}

func Status(status string) Option {
	return func(b *Job) error {
		_, ok := JOB_STATUS[status]
		if !ok {
			return fmt.Errorf("unsupported status [%s]", status)
		}
		b.Status = status
		return nil
	}
}
