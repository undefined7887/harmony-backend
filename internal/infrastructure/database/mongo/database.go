package mongodatabase

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/undefined7887/harmony-backend/internal/config"
)

func NewDatabase(config *config.Mongo) (*mongo.Database, error) {
	clientOpts := options.Client().
		SetDirect(config.Direct).
		ApplyURI(fmt.Sprintf("mongodb://%s", config.Address))

	if config.Username != "" && config.Password != "" {
		clientOpts = clientOpts.SetAuth(options.Credential{
			AuthMechanism: "PLAIN",
			Username:      config.Username,
			Password:      config.Password,
		})
	}

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}

	return client.Database(config.Database), nil
}

func NewDatabaseRunner(
	lifecycle fx.Lifecycle,
	config *config.Mongo,
	logger *zap.Logger,
	database *mongo.Database,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("connecting to mongo database", zap.String("address", config.Address))

			client := database.Client()

			if err := client.Connect(ctx); err != nil {
				return fmt.Errorf("mongo: %v", err)
			}

			if err := client.Ping(ctx, nil); err != nil {
				return fmt.Errorf("mongo: %v", err)
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("disconnecting from mongo database")

			return database.Client().Disconnect(ctx)
		},
	})
}

func Transaction[T any](
	ctx context.Context,
	database *mongo.Database,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var emptyResult T

	session, err := database.Client().StartSession()
	if err != nil {
		return emptyResult, err
	}

	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (any, error) {
		return fn(sessionCtx)
	})
	if err != nil {
		return emptyResult, err
	}

	return result.(T), nil
}

func TransactionNoReturn(
	ctx context.Context,
	database *mongo.Database,
	fn func(ctx context.Context) error,
) error {
	_, err := Transaction[struct{}](ctx, database, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(ctx)
	})

	return err
}

func IndexKeys(keys ...string) bson.D {
	return keysToMap(keys, 1)
}

func IndexKeysCustom(value any, keys ...string) bson.D {
	return keysToMap(keys, value)
}

func keysToMap(keys []string, value any) bson.D {
	result := make(bson.D, len(keys))

	for i, key := range keys {
		result[i] = bson.E{Key: key, Value: value}
	}

	return result
}
