package words

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(myContext context.Context, word *Word) (string, error)
}

type repo struct {
	collection *mongo.Collection
}

type collectionGetter interface {
	Collection(string) *mongo.Collection
}

func NewRepository(dataBase collectionGetter) Repository {
	return &repo{collection: dataBase.Collection("words")}
}

func (repository *repo) Create(myContext context.Context, word *Word) (string, error) {
	now := time.Now().UTC()
	word.CreatedAt = now

	// We let Mongo assign _id; we just return it as hex
	response, err := repository.collection.InsertOne(myContext, bson.M{
		"text":      word.Text,
		"ownerId":   word.OwnerID,
		"createdAt": word.CreatedAt,
	})
	if err != nil {
		return "", err
	}
	oid, ok := response.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNilDocument
	}
	return oid.Hex(), nil
}
