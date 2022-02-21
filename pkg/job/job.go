package job

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var JobRequests = make(chan Job)
var validate *validator.Validate

type Job struct {
	ID          string `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string `bson:"type" json:"type" validate:"required"`
	Status      string `bson:"status" json:"status" validate:"required"`
	InitiatedBy string `bson:"initiated_by" json:"initiated_by" validate:"required"`
	SourceFile  string `bson:"source_file" json:"source_file"`
	Created     string `bson:"created" json:"created"`
	Modified    string `bson:"modified" json:"modified"`
}

func NewJob(opts ...Option) (*Job, error) {
	u := &Job{
		Status: JOB_STATUS_PENDING,
	}
	for _, o := range opts {
		if err := o(u); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (b *Job) Validate() error {
	validate = validator.New()
	err := validate.Struct(b)
	if err != nil {
		return err
	}
	return nil
}

func (b *Job) GetObjectID() primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(b.ID)
	if err == nil {
		return oid
	}
	return primitive.NilObjectID
}

func (b *Job) SetStatus(status string) error {
	opts := Status(status)
	return opts(b)
}

func CreateNewJob(c context.Context, jobType string, sourceFile string) (*Job, error) {
	batch, err := NewJob(
		Type(jobType),
		Status(JOB_STATUS_PENDING),
		InitiatedBy(utils.GetUserEmail(c)),
		SourceFile(sourceFile),
	)
	if err != nil {
		return nil, err
	}

	err = batch.Store(c)
	if err != nil {
		return nil, err
	}
	go AddToJobQueue(*batch)
	return batch, err
}
