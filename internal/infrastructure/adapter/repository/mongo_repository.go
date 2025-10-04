package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"telecomx-support-service/internal/domain/model"
)

type MongoRepository struct {
	Collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{Collection: db.Collection("support")}
}

func (r *MongoRepository) Create(ctx context.Context, p *model.Support) error {
	p.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, p)
	return err
}

func (r *MongoRepository) UpdateStatus(ctx context.Context, userID, status string) error {
	filter := bson.M{"userId": userID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"userId": userID})
	return err
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]model.Support, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Support
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
