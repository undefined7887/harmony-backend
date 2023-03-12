package mongodatabase

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/samber/lo"
)

const (
	batchSize = 1000
)

type Query[T any] struct {
	collection *mongo.Collection
}

func NewQuery[T any](collection *mongo.Collection) *Query[T] {
	return &Query[T]{
		collection: collection,
	}
}

func (q *Query[T]) BuildIndex(ctx context.Context, keys interface{}, opts ...*options.IndexOptions) error {
	_, err := q.collection.
		Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys:    keys,
			Options: options.MergeIndexOptions(opts...),
		})

	return err
}

func (q *Query[T]) InsertOne(ctx context.Context, model *T) (bool, error) {
	_, err := q.collection.InsertOne(ctx, model)

	if mongo.IsDuplicateKeyError(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (q *Query[T]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error) {
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

func (q *Query[T]) FindCursor(
	ctx context.Context,
	filter bson.M,
	cb func(ctx context.Context, models []T) error,
	opts ...*options.FindOptions,
) error {
	cursor, err := q.collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	batch := make([]T, 0, batchSize)

	for cursor.Next(ctx) {
		var result T

		if err := cursor.Decode(&result); err != nil {
			return err
		}

		batch = append(batch, result)

		if len(batch) >= batchSize {
			if err := cb(ctx, batch); err != nil {
				return err
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := cb(ctx, batch); err != nil {
			return err
		}
	}

	if cursor.Err() != nil {
		return err
	}

	return nil
}

func (q *Query[T]) FindOne(ctx context.Context, filter bson.M, opts ...*options.FindOneOptions) (T, error) {
	var result T

	err := q.collection.
		FindOne(ctx, filter, opts...).
		Decode(&result)

	if err != nil {
		return lo.Empty[T](), err
	}

	return result, nil
}

func (q *Query[T]) Exists(ctx context.Context, filter bson.M) (bool, error) {
	result, err := q.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

func (q *Query[T]) Aggregate(ctx context.Context, pipeline bson.A, opts ...*options.AggregateOptions) ([]T, error) {
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

func (q *Query[T]) UpdateMany(ctx context.Context, filter bson.M, update any, opts ...*options.UpdateOptions) (int64, error) {
	result, err := q.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (q *Query[T]) FindOneAndUpdate(ctx context.Context, filter bson.M, update any, opts ...*options.FindOneAndUpdateOptions) (T, error) {
	var result T

	err := q.collection.
		FindOneAndUpdate(ctx, filter, update, opts...).
		Decode(&result)

	if err != nil {
		return lo.Empty[T](), err
	}

	return result, nil
}
