package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func (myMongo *Mongo) People() *mongo.Collection {
	return myMongo.DB.Collection("people")
}

func (myMongo *Mongo) Close(ctx context.Context) error {
	return myMongo.Client.Disconnect(ctx)
}

func Connect(ctx context.Context, uri, dbName string) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := c.Ping(ctx, nil); err != nil {
		return nil, err
	}

	myMongo := &Mongo{Client: c, DB: c.Database(dbName)}

	// Ensure a unique index on email
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_email"),
	}
	_, _ = myMongo.People().Indexes().CreateOne(ctx, index)

	return myMongo, nil
}
