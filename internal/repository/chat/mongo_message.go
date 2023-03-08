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

type MongoMessageRepository struct {
	database *mongo.Database
}

func NewMongoMessageRepository(database *mongo.Database) chatdomain.MessageRepository {
	return &MongoMessageRepository{
		database: database,
	}
}

func NewMongoMessageMigrationsRunner(lifecycle fx.Lifecycle, logger *zap.Logger, database *mongo.Database) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running migrations", zap.String("collection", messageCollection))

			return multierr.Combine(
				mongodatabase.
					NewQuery[any](database.Collection(messageCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("chat_id"),
					),

				mongodatabase.
					NewQuery[any](database.Collection(messageCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("user_id", "peer_id"),
					),

				mongodatabase.
					NewQuery[any](database.Collection(messageCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("user_read_ids"),
					),
			)
		},
	})
}

func (m *MongoMessageRepository) Create(ctx context.Context, message *chatdomain.Message) (bool, error) {
	return mongodatabase.
		NewQuery[chatdomain.Message](m.database.Collection(messageCollection)).
		InsertOne(ctx, message)
}

func (m *MongoMessageRepository) Get(ctx context.Context, id string) (*chatdomain.Message, error) {
	return mongodatabase.
		NewQuery[chatdomain.Message](m.database.Collection(messageCollection)).
		FindOne(ctx, bson.M{
			"_id": id,
		})
}

func (m *MongoMessageRepository) List(ctx context.Context, chatID string, offset, limit int64) ([]chatdomain.Message, error) {
	return mongodatabase.
		NewQuery[chatdomain.Message](m.database.Collection(messageCollection)).
		Find(ctx,
			bson.M{
				"chat_id": chatID,
			},
			options.Find().
				SetSort(bson.M{
					"created_at": -1,
				}).
				SetSkip(offset).
				SetLimit(limit),
		)
}

func (m *MongoMessageRepository) UpdateText(ctx context.Context, id, userID, text string) (*chatdomain.Message, error) {
	return mongodatabase.
		NewQuery[chatdomain.Message](m.database.Collection(messageCollection)).
		FindOneAndUpdate(ctx,
			bson.M{
				"_id": id,

				// Only sender can modify message
				"user_id": userID,
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
