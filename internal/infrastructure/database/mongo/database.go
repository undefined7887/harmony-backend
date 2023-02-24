package mongodatabase

import (
	"context"
	"fmt"
	"github.com/undefined7887/harmony-backend/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
				return err
			}

			return client.Ping(ctx, nil)
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
