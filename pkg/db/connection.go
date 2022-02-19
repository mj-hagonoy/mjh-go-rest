package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	once   sync.Once
	client *mongo.Client
}

var database Database

func ConnectDB() (*mongo.Client, error) {
	database.once.Do(createInstance)

	return database.client, nil
}

func createInstance() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%d", config.GetConfig().Database.Host, config.GetConfig().Database.Port)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	database.client = mongoClient
}
