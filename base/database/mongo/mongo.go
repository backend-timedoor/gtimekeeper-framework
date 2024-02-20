package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(url string) (*mongo.Client, error) {
	mongo, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return mongo, nil
}
