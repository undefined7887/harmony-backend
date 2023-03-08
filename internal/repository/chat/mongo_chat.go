package chatrepo

import (
	"context"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	mongodatabase "github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	chatCollection = "chats"
)

type MongoChatRepository struct {
	database *mongo.Database
}

func NewMongoChatRepository(db *mongo.Database) chatdomain.ChatRepository {
	return &MongoChatRepository{
		database: db,
	}
}

func (m *MongoChatRepository) List(ctx context.Context, userID, chatType string, offset, limit int64) ([]chatdomain.Chat, error) {
	match := bson.M{
		"$or": bson.A{
			bson.M{"user_id": userID},
			bson.M{"peer_id": userID},
		},
	}

	// Add filter on chat_type if presented
	if chatType != "" {
		match["chat_type"] = chatType
	}

	return mongodatabase.
		NewQuery[chatdomain.Chat](m.database.Collection(messageCollection)).
		Aggregate(ctx, bson.A{
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
					"from": chatCollection,
					"as":   "chat",

					"localField":   "message.chat_id",
					"foreignField": "_id",
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$chat",

					// For more information, see:
					// https://www.mongodb.com/docs/manual/reference/operator/aggregation/unwind
					"preserveNullAndEmptyArrays": true,
				},
			},
			bson.M{
				"$group": bson.M{
					"_id": "$message.chat_id",
					"chat": bson.M{
						"$first": "$chat",
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
											"$setIsSubset": bson.A{bson.A{userID}, "$message.user_read_ids"},
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
							"$chat",
							bson.M{
								// If chat not found
								"_id":  "$message.chat_id",
								"type": "$message.chat_type",

								"message": bson.M{
									"last":         "$message",
									"unread_count": "$unread_count",
								},
							},
						},
					},
				},
			},
			bson.M{
				"$replaceRoot": bson.M{
					"newRoot": "$chat",
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

func (m *MongoChatRepository) UpdateRead(ctx context.Context, userID, chatID string) (int64, error) {
	return mongodatabase.
		NewQuery[chatdomain.Message](m.database.Collection(messageCollection)).
		UpdateMany(ctx,
			bson.M{
				"user_id": bson.M{"$ne": userID},
				"chat_id": chatID,
				"user_read_ids": bson.M{
					"$nin": bson.A{userID},
				},
			},
			bson.M{
				"$addToSet": bson.M{
					"user_read_ids": userID,
				},
				"$set": bson.M{
					"updated_at": time.Now(),
				},
			},
		)
}
