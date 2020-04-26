package db

import (
	"context"
	"fmt"
	"github.com/ToolPackage/fse/server/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	MongoClient = NewMongoClient()
	MongoDb     = MongoClient.Database(config.MongoDbName)
)

func NewMongoClient() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	uri := fmt.Sprintf("mongodb://%s:%s", config.Conf.MongoHost, config.Conf.MongoPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(errors.WithStack(err))
	}
	return client
}
