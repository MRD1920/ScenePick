package DBConfig

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfigMgr struct {
	MongoClient *mongo.Client
	ctx         context.Context
	cancel      context.CancelFunc
}

func ConnectToMongoDB(MongoConnectionURI string) (*DBConfigMgr, error) {
	//Get the connection string from config of server
	uri := MongoConnectionURI
	if uri == "" {
		log.Fatal("MongoDB connection string is not provided")
		return nil, fmt.Errorf("MongoDB connection string is not provided")
	}

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("Cannot connect to MongoDB: ", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Cannot ping to MongoDB: ", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB !!")
	return &DBConfigMgr{MongoClient: client, ctx: ctx, cancel: cancel}, nil

}
