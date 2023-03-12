package userrepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/undefined7887/harmony-backend/internal/domain/user"
	"github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	"github.com/undefined7887/harmony-backend/internal/util"
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
					NewQuery[any](database.Collection(userCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("email"),
						options.
							Index().
							SetUnique(true),
					),
			)
		},
	})
}

func (m *MongoRepository) Create(ctx context.Context, user *userdomain.User) (bool, error) {
	return mongodatabase.
		NewQuery[userdomain.User](m.database.Collection(userCollection)).
		InsertOne(ctx, user)
}

func (m *MongoRepository) Get(ctx context.Context, id string) (userdomain.User, error) {
	return mongodatabase.
		NewQuery[userdomain.User](m.database.Collection(userCollection)).
		FindOne(ctx, bson.M{
			"_id": id,
		})
}

func (m *MongoRepository) GetByEmail(ctx context.Context, email string) (userdomain.User, error) {
	return mongodatabase.
		NewQuery[userdomain.User](m.database.Collection(userCollection)).
		FindOne(ctx, bson.M{
			"email": email,
		})
}

func (m *MongoRepository) GetByNickname(ctx context.Context, nickname string) (userdomain.User, error) {
	return mongodatabase.
		NewQuery[userdomain.User](m.database.Collection(userCollection)).
		FindOne(ctx, bson.M{
			"nickname": nickname,
		})
}

func (m *MongoRepository) Exists(ctx context.Context, id string) (bool, error) {
	return mongodatabase.
		NewQuery[userdomain.User](m.database.Collection(userCollection)).
		Exists(ctx, bson.M{
			"_id": id,
		})
}

func (m *MongoRepository) UpdateStatus(ctx context.Context, id, status string, onlyOffline bool) (userdomain.User, error) {
	var statusSet any = status

	if onlyOffline {
		// If 'onlyOffline' flag is set change status only if it equals to 'offline'
		statusSet = bson.M{
			"$cond": bson.M{
				"if": bson.M{
					"$eq": bson.A{"$status", userdomain.StatusOffline},
				},
				"then": status,
				"else": "$status",
			},
		}
	}

	return mongodatabase.
		NewQuery[userdomain.User](m.database.Collection(userCollection)).
		FindOneAndUpdate(ctx,
			bson.M{
				"_id": id,
			},
			bson.A{
				bson.M{
					"$set": bson.M{
						"status":     statusSet,
						"updated_at": time.Now(),
					},
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.After),
		)
}

func (m *MongoRepository) UpdateOutdatedStatuses(ctx context.Context, cb func(users []userdomain.User)) error {
	collection := m.database.Collection(userCollection)

	now := time.Now()

	batchFunc := func(ctx context.Context, users []userdomain.User) error {
		_, err := collection.BulkWrite(ctx, util.Map(users, func(user userdomain.User) mongo.WriteModel {
			return &mongo.UpdateOneModel{
				Filter: bson.M{
					"_id": user.ID,
				},
				Update: bson.M{
					"$set": bson.M{
						"status":     userdomain.StatusOffline,
						"updated_at": now,
					},
				},
			}
		}))

		if err != nil {
			return err
		}

		// Before calling callback function we need to update users
		users = util.Map(users, func(user userdomain.User) userdomain.User {
			user.Status = userdomain.StatusOffline
			user.UpdatedAt = now

			return user
		})

		// Some users can be processed even if transaction failed
		cb(users)

		return nil
	}

	return mongodatabase.TransactionNoReturn(ctx, m.database, func(ctx context.Context) error {
		return mongodatabase.
			NewQuery[userdomain.User](collection).
			FindCursor(ctx,
				bson.M{
					"status": bson.M{
						"$ne": userdomain.StatusOffline,
					},
					"updated_at": bson.M{
						"$lt": time.Now().Add(-userdomain.UserOutdatedTimeout),
					},
				},
				batchFunc,
			)
	})
}
