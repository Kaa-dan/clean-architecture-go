package repositories

import (
	"context"
	"time"

	"github.com/kaa-dan/clean-architecture-go/internal/domain/enitities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *userRepository) Create(ctx context.Context, user *enitities.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrUserAlredyExitst
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil

}
