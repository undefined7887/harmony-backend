package callrepo

import (
	"context"
	"encoding/json"
	"fmt"
	calldomain "github.com/undefined7887/harmony-backend/internal/domain/call"
	mongodatabase "github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	"github.com/undefined7887/harmony-backend/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	callCollection = "calls"
)

type MongoRepository struct {
	database *mongo.Database
}

func NewMongoRepository(database *mongo.Database) calldomain.Repository {
	return &MongoRepository{
		database: database,
	}
}

func NewMongoMigrationsRunner(lifecycle fx.Lifecycle, logger *zap.Logger, database *mongo.Database) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running migrations", zap.String("collection", callCollection))

			return multierr.Combine(
				mongodatabase.
					NewQuery[any](database.Collection(callCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("user_id"),
					),

				mongodatabase.
					NewQuery[any](database.Collection(callCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("peer_id"),
					),

				mongodatabase.
					NewQuery[any](database.Collection(callCollection)).
					BuildIndex(ctx,
						mongodatabase.IndexKeys("status"),
					),
			)
		},
	})
}

func (m *MongoRepository) Create(ctx context.Context, call *calldomain.Call) (bool, error) {
	result, err := m.database.
		Collection(callCollection).
		UpdateOne(ctx,
			// Checking that potential call members don't have current active requests
			bson.M{
				"$or": bson.A{
					bson.M{"user_id": call.UserID},
					bson.M{"user_id": call.PeerID},
					bson.M{"peer_id": call.UserID},
					bson.M{"peer_id": call.PeerID},
				},
				"status": calldomain.StatusRequest,
			},
			bson.M{
				"$setOnInsert": call,
			},
			options.
				Update().
				SetUpsert(true),
		)

	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, nil
}

func (m *MongoRepository) Read(ctx context.Context, id, status string) (calldomain.Call, error) {
	return mongodatabase.
		NewQuery[calldomain.Call](m.database.Collection(callCollection)).
		FindOne(ctx, bson.M{
			"_id":    id,
			"status": status,
		})
}

func (m *MongoRepository) ReadLast(ctx context.Context, userID, status string) (calldomain.Call, error) {
	return mongodatabase.
		NewQuery[calldomain.Call](m.database.Collection(callCollection)).
		FindOne(ctx, bson.M{
			"$or": bson.A{
				bson.M{"user_id": userID},
				bson.M{"peer_id": userID},
			},
			"status": status,
		})
}

func (m *MongoRepository) UpdateStatus(ctx context.Context, id, peerID string, previousStatuses []string, newStatus string) (calldomain.Call, error) {
	match := bson.M{
		"_id": id,
		"status": bson.M{
			// Converting to []any
			"$in": util.Map(previousStatuses, func(item string) any {
				return item
			}),
		},
	}

	if peerID != "" {
		match["peer_id"] = peerID
	}

	res, _ := json.Marshal(match)
	fmt.Println(string(res))

	return mongodatabase.NewQuery[calldomain.Call](m.database.Collection(callCollection)).
		FindOneAndUpdate(ctx,
			match,
			bson.M{
				"$set": bson.M{
					"status": newStatus,
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.After),
		)
}
