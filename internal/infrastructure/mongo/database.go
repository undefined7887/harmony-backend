package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/undefined7887/harmony-backend/internal/config"
)

type Database struct {
	*mongo.Database
}

func NewDatabase(ctx context.Context, config *config.Mongo) (*Database, error) {
	client, err := mongo.Connect(ctx,
		options.
			Client().
			SetAuth(options.Credential{
				AuthMechanism: "PLAIN",
				Username:      config.Username,
				Password:      config.Password,
			}).
			ApplyURI(fmt.Sprintf("mongodb://%s", config.Address)),
	)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Database{
		Database: client.Database(config.Database),
	}, nil
}

type TransactionFunc func(ctx context.Context) (any, error)

func (d *Database) Transaction(ctx context.Context, fn TransactionFunc) (any, error) {
	session, err := d.Client().StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (any, error) {
		return fn(sessionCtx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
