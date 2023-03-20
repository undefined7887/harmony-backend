package chatrepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	mongodatabase "github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
)

type MongoChatRepository struct {
	database *mongo.Database
}

func NewMongoChatRepository(db *mongo.Database) chatdomain.ChatRepository {
	return &MongoChatRepository{
		database: db,
	}
}

func (m *MongoChatRepository) List(ctx context.Context, userID, peerType string, offset, limit int64) ([]chatdomain.Chat, error) {
	match := bson.M{
		"$or": bson.A{
			bson.M{"user_id": userID},
			bson.M{"peer_id": userID},
		},
	}

	// Add filter on chat_type if presented
	if peerType != "" {
		match["peer_type"] = peerType
	}

	pipeline := bson.A{
		bson.M{
			"$match": match,
		},
		bson.M{
			"$sort": bson.M{
				"created_at": -1,
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"message": "$$ROOT",
				},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from": groupCollection,
				"as":   "group",

				"localField":   "message.chat_id",
				"foreignField": "_id",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$group",

				// For more information, see:
				// https://www.mongodb.com/docs/manual/reference/operator/aggregation/unwind
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$message.chat_id",
				"group": bson.M{
					"$first": "$group",
				},
				"message": bson.M{
					"$first": "$message",
				},
				"unread_count": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							// Do not count if:
							// - current user is sender
							// or
							// - current user read message
							"if": bson.M{
								"$or": bson.A{
									bson.M{
										"$eq": bson.A{userID, "$message.user_id"},
									},
									bson.M{
										"$setIsSubset": bson.A{bson.A{userID}, "$message.read_user_ids"},
									},
								},
							},
							"then": 0,
							"else": 1,
						},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"chat": bson.M{
					"$mergeObjects": bson.A{
						bson.M{
							"message":      "$message",
							"unread_count": "$unread_count",
						},
						// User chat
						bson.M{
							// For user chats choosing chat_id as opposite to current user
							"_id": bson.M{
								"$cond": bson.M{
									"if": bson.M{
										"$eq": bson.A{userID, "$message.user_id"},
									},
									"then": "$message.peer_id",
									"else": "$message.user_id",
								},
							},
							"type": "$message.peer_type",
						},

						// Group chat
						"$group",
					},
				},
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$chat",
			},
		},
	}

	if offset > 0 {
		pipeline = append(pipeline, bson.M{
			"$skip": offset,
		})
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.M{
			"$limit": limit,
		})
	}

	return mongodatabase.
		NewQuery[chatdomain.Chat](m.database.Collection(messageCollection)).
		Aggregate(ctx, pipeline)
}

func (m *MongoChatRepository) UpdateRead(ctx context.Context, userID, chatID string) (int64, error) {
	return mongodatabase.
		NewQuery[chatdomain.Message](m.database.Collection(messageCollection)).
		UpdateMany(ctx,
			bson.M{
				"chat_id": chatID,

				// User can't read self messages
				"user_id": bson.M{"$ne": userID},

				// User can't re-read message
				"read_user_ids": bson.M{
					"$nin": bson.A{userID},
				},
			},
			bson.M{
				"$addToSet": bson.M{
					"read_user_ids": userID,
				},
				"$set": bson.M{
					"updated_at": time.Now(),
				},
			},
		)
}
