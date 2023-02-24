package mongodatabase

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionQuery[T any] struct {
	collection *mongo.Collection
}

func NewCollectionQuery[T any](collection *mongo.Collection) *CollectionQuery[T] {
	return &CollectionQuery[T]{
		collection: collection,
	}
}
func (q *CollectionQuery[T]) BuildIndex(ctx context.Context, keys interface{}, options *options.IndexOptions) error {
	_, err := q.collection.
		Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys:    keys,
			Options: options,
		})

	return err
}

func (q *CollectionQuery[T]) InsertOne(ctx context.Context, model *T) (bool, error) {
	_, err := q.collection.InsertOne(ctx, model)

	if mongo.IsDuplicateKeyError(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (q *CollectionQuery[T]) FindOne(ctx context.Context, filter bson.M) (*T, error) {
	var result T

	err := q.collection.
		FindOne(ctx, filter).
		Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
