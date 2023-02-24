package userrepo

import (
	"context"
	"github.com/undefined7887/harmony-backend/internal/domain/user"
	"github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	userCollection = "users"
)

type MongoRepository struct {
	database *mongo.Database
}

func NewMongoRepository(database *mongo.Database) userdomain.Repository {
	return &MongoRepository{
		database: database,
	}
}

func NewMongoMigrationsRunner(lifecycle fx.Lifecycle, logger *zap.Logger, database *mongo.Database) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running migrations", zap.String("collection", userCollection))

			return multierr.Combine(
				mongodatabase.
					NewCollectionQuery[any](database.Collection(userCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("email"),
						options.Index().
							SetUnique(true),
					),
			)
		},
	})
}

func (m *MongoRepository) Create(ctx context.Context, user *userdomain.User) (bool, error) {
	return mongodatabase.
		NewCollectionQuery[userdomain.User](m.database.Collection(userCollection)).
		InsertOne(ctx, user)
}

func (m *MongoRepository) Read(ctx context.Context, id string) (*userdomain.User, error) {
	return mongodatabase.
		NewCollectionQuery[userdomain.User](m.database.Collection(userCollection)).
		FindOne(ctx, bson.M{
			"id": id,
		})
}

func (m *MongoRepository) ReadByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	return mongodatabase.
		NewCollectionQuery[userdomain.User](m.database.Collection(userCollection)).
		FindOne(ctx, bson.M{
			"email": email,
		})
}
