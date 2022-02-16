package job

import (
	"context"
	"fmt"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (job *Job) Store(c context.Context) error {
	dbc, err := db.ConnectDB()
	if err != nil {
		return err
	}
	coll := dbc.Database(config.GetConfig().Database.DbName).Collection(db.COL_JOBS)
	result, err := coll.InsertOne(c, job)
	if err != nil {
		return err
	}
	job.ID = (result.InsertedID.(primitive.ObjectID)).Hex()
	return nil
}

func (job *Job) GetOne(c context.Context) error {
	dbc, err := db.ConnectDB()
	if err != nil {
		return err
	}
	coll := dbc.Database(config.GetConfig().Database.DbName).Collection(db.COL_JOBS)
	filter := bson.M{"_id": job.GetObjectID()}

	result := coll.FindOne(c, filter)
	if result.Err() != nil {
		return fmt.Errorf("error retrieving record [%s]", result.Err().Error())
	}
	if err := result.Decode(job); err != nil {
		return err
	}
	return nil
}

func (job *Job) Update(c context.Context) error {
	dbc, err := db.ConnectDB()
	if err != nil {
		return err
	}
	coll := dbc.Database(config.GetConfig().Database.DbName).Collection(db.COL_JOBS)
	filter := bson.M{"_id": job.GetObjectID()}
	update := bson.M{
		"$set": bson.M{
			"status": job.Status,
		},
	}

	_, err = coll.UpdateOne(c, filter, update)
	if err != nil {
		return fmt.Errorf("error updating record [%s]", err.Error())
	}
	return nil
}
