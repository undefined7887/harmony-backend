package repository

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

func IsNoDocumentsErr(err error) bool {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return true
	}

	return false
}
