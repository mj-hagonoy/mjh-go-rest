package user

import (
	"context"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (user *User) Store(c context.Context) (string, error) {
	dbc, err := db.ConnectDB()
	if err != nil {
		return "", err
	}
	coll := dbc.Database(config.GetConfig().Database.DbName).Collection(db.COL_USERS)
	result, err := coll.InsertOne(c, user)
	if err != nil {
		return "", err
	}
	return (result.InsertedID.(primitive.ObjectID)).Hex(), nil
}

func GetAll(c context.Context) (Users, error) {
	dbc, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}
	coll := dbc.Database(config.GetConfig().Database.DbName).Collection(db.COL_USERS)
	cursor, err := coll.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}
	var records Users
	if err := cursor.All(c, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (users Users) StoreBulk(c context.Context) (int, error) {
	dbc, err := db.ConnectDB()
	if err != nil {
		return 0, err
	}
	coll := dbc.Database(config.GetConfig().Database.DbName).Collection(db.COL_USERS)
	result, err := coll.InsertMany(c, toInterface(users))
	if err != nil {
		return 0, err
	}
	return len(result.InsertedIDs), nil
}

func toInterface(users Users) []interface{} {
	y := make([]interface{}, len(users))
	for i, v := range users {
		y[i] = v
	}
	return y
}
