package chatrepo

import (
	"context"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	mongodatabase "github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"time"
)

const (
	messageCollection = "messages"
)

type MongoRepository struct {
	database *mongo.Database
}

func NewMongoRepository(database *mongo.Database) chatdomain.Repository {
	return &MongoRepository{
		database: database,
	}
}

func NewMongoMigrationsRunner(lifecycle fx.Lifecycle, logger *zap.Logger, database *mongo.Database) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running migrations", zap.String("collection", messageCollection))

			return multierr.Combine(
				mongodatabase.
					Q[any](database.Collection(messageCollection)).
					BuildIndex(ctx, mongodatabase.IndexKeys("user_id")),

				mongodatabase.
					Q[any](database.Collection(messageCollection)).
					BuildIndex(ctx, mongodatabase.IndexKeys("peer_id")),

				mongodatabase.
					Q[any](database.Collection(messageCollection)).
					BuildIndex(ctx, mongodatabase.IndexKeys("peer_hash")),
			)
		},
	})
}

func (m *MongoRepository) Create(ctx context.Context, message *chatdomain.Message) (bool, error) {
	return mongodatabase.
		Q[chatdomain.Message](m.database.Collection(messageCollection)).
		InsertOne(ctx, message)
}

func (m *MongoRepository) List(ctx context.Context, peerHash string, offset, limit int64) ([]chatdomain.Message, error) {
	return mongodatabase.
		Q[chatdomain.Message](m.database.Collection(messageCollection)).
		Find(ctx,
			bson.M{
				"peer_hash": peerHash,

				// Do not include deleted messages
				"deleted_at": nil,
			},
			options.Find().
				SetSort(bson.M{"created_at": -1}).
				SetSkip(offset).
				SetLimit(limit),
		)
}

func (m *MongoRepository) ListRecent(ctx context.Context, userID, peerType string, offset, limit int64) ([]chatdomain.Message, error) {
	match := bson.M{
		"$or": bson.A{
			bson.M{"user_id": userID},
			bson.M{"peer_id": userID},
		},

		// Do not include deleted messages
		"deleted_at": nil,
	}

	// Add filter on peer_type if presented
	if peerType != "" {
		match["peer_type"] = peerType
	}

	return mongodatabase.
		Q[chatdomain.Message](m.database.Collection(messageCollection)).
		Aggregate(ctx, bson.A{
			// In future this request will $lookup groups
			bson.M{
				"$match": match,
			},
			bson.M{
				"$sort": bson.M{
					"created_at": -1,
				},
			},
			bson.M{
				"$group": bson.M{
					"_id":  "$peer_hash",
					"root": bson.M{"$first": "$$ROOT"},
				},
			},
			bson.M{
				"$replaceRoot": bson.M{
					"newRoot": "$root",
				},
			},
			bson.M{
				"$skip": offset,
			},
			bson.M{
				"$limit": limit,
			},
		})
}

func (m *MongoRepository) Update(ctx context.Context, id, userID, text string) (*chatdomain.Message, error) {
	return mongodatabase.
		Q[chatdomain.Message](m.database.Collection(messageCollection)).
		FindOneAndUpdate(ctx,
			bson.M{
				"_id":     id,
				"user_id": userID,

				// Do not include deleted messages
				"deleted_at": nil,
			},
			bson.M{
				"$set": bson.M{
					"text":       text,
					"updated_at": time.Now(),
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.After),
		)
}

func (m *MongoRepository) UpdateRead(ctx context.Context, userID, peerHash string) (int64, error) {
	return mongodatabase.
		Q[chatdomain.Message](m.database.Collection(messageCollection)).
		UpdateMany(ctx,
			bson.M{
				"user_id":   userID,
				"peer_hash": peerHash,

				// Do not include deleted messages
				"deleted_at": nil,
			},
			bson.M{
				"$set": bson.M{
					"read":       true,
					"updated_at": time.Now(),
				},
			},
		)
}
