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

func Q[T any](collection *mongo.Collection) *CollectionQuery[T] {
	return &CollectionQuery[T]{
		collection: collection,
	}
}
func (q *CollectionQuery[T]) BuildIndex(ctx context.Context, keys interface{}, opts ...*options.IndexOptions) error {
	_, err := q.collection.
		Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys:    keys,
			Options: options.MergeIndexOptions(opts...),
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

func (q *CollectionQuery[T]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := q.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var result []T
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (q *CollectionQuery[T]) FindOne(ctx context.Context, filter bson.M, opts ...*options.FindOneOptions) (*T, error) {
	var result T

	err := q.collection.
		FindOne(ctx, filter, opts...).
		Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (q *CollectionQuery[T]) Exists(ctx context.Context, filter bson.M) (bool, error) {
	result, err := q.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

func (q *CollectionQuery[T]) Aggregate(ctx context.Context, pipeline bson.A, opts ...*options.AggregateOptions) ([]T, error) {
	cursor, err := q.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	var result []T
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (q *CollectionQuery[T]) UpdateMany(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, error) {
	result, err := q.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (q *CollectionQuery[T]) FindOneAndUpdate(ctx context.Context, filter bson.M, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	var result T

	err := q.collection.
		FindOneAndUpdate(ctx, filter, update, opts...).
		Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
