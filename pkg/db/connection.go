package db

import (
	"context"
	"fmt"
	"time"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDB() (*mongo.Client, error) {
	if mongoClient == nil {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		uri := fmt.Sprintf("mongodb://%s:%d", config.GetConfig().Database.Host, config.GetConfig().Database.Port)
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
	}

	return mongoClient, nil
}
