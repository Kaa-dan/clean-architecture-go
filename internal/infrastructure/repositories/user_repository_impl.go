package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) *userRepository {
	collection := client.Database(dbName).Collection("users")

	//Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	//Email index (unique)

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	//Username index (unique)
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	collection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, usernameIndex})

	return &userRepository{
		collection: collection,
	}
}
