package repository

import (
	"context"
	"github.com/sgitwhyd/jagong/app/models"
	"github.com/sgitwhyd/jagong/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertMessage(ctx context.Context, model models.MessagePayload) error {
	_, err := database.MongoDB.InsertOne(ctx, model)
	return err
}

func FindAllMessage(ctx context.Context) ([]models.MessagePayload, error) {
	filter := bson.M{}
	cursor, err := database.MongoDB.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result []models.MessagePayload
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
