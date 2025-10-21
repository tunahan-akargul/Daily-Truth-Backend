package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Collection is a generic method to get any collection by name
func (myMongo *Mongo) Collection(name string) *mongo.Collection {
	return myMongo.DB.Collection(name)
}

func (myMongo *Mongo) Close(ctx context.Context) error {
	return myMongo.Client.Disconnect(ctx)
}

func Connect(ctx context.Context, uri, dbName string) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	myMongo := &Mongo{Client: client, DB: client.Database(dbName)}

	return myMongo, nil
}
