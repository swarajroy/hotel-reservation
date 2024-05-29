package utils

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func ResourceNotFound(name, id string, err error) error {
	if errors.Is(mongo.ErrNoDocuments, err) {
		return fmt.Errorf("resource %s with id %s not found", name, id)
	}
	return err
}
