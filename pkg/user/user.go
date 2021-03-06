package user

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/utils"
)

var validate *validator.Validate

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"first_name" json:"first_name" validate:"required"`
	LastName  string `bson:"last_name" json:"last_name" validate:"required"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"-"`
	Activated bool   `bson:"activated" json:"activated"`
	Created   string `bson:"created" json:"created"`
	Modified  string `bson:"modified" json:"modified"`
}
type Users []*User

func NewUser(opts ...Option) (*User, error) {
	u := &User{}
	for _, o := range opts {
		if err := o(u); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (u *User) Validate() error {
	validate = validator.New()
	err := validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Bind(data []string) error {
	if len(data) < 3 {
		return fmt.Errorf("size must be 3")
	}

	u.FirstName = data[0]
	u.LastName = data[1]
	u.Email = data[2]
	return nil
}

func ImportUsersFromCsv(ctx context.Context, filepath string) error {
	records, err := utils.CsvReadAll(filepath, true)
	if err != nil {
		return err
	}
	return BulkUpload(ctx, records)
}

func ImportUsersFromCsvBytes(ctx context.Context, data []byte) error {
	buf := bytes.NewReader(data)
	records, err := utils.CsvRead(buf, true)
	if err != nil {
		return fmt.Errorf("ImportUsersFromCsvBytes: %v", err)
	}

	return BulkUpload(ctx, records)
}

func BulkUpload(ctx context.Context, records [][]string) error {
	if len(records) == 0 {
		return nil
	}
	batchSize, size := 500, len(records)
	var wg sync.WaitGroup
	for i := 0; i < size; {
		start, end := i, i+batchSize
		if end > size {
			end = size
		}
		wg.Add(1)
		go bulkUpload(ctx, &wg, records[start:end])
		i += batchSize
	}
	wg.Wait()
	return nil
}

func bulkUpload(ctx context.Context, wg *sync.WaitGroup, records [][]string) {
	if len(records) == 0 {
		return
	}
	var toInsertUsers Users = make(Users, 0)
	now := time.Now().UTC().String()

	for _, record := range records {
		user, e := NewUser()
		if e != nil {
			logger.ErrorLogger.Printf("user instance creation issue with error: [%s]", e.Error())
			continue
		}
		if e := user.Bind(record); e != nil {
			logger.ErrorLogger.Printf("bind failed with error: [%s]", e.Error())
			continue
		}
		user.Created = now
		user.Modified = now
		toInsertUsers = append(toInsertUsers, user)
	}
	_, err := toInsertUsers.StoreBulk(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		logger.ErrorLogger.Printf("bulk insert failed: [%s]", err.Error())
	}
	wg.Done()
}
