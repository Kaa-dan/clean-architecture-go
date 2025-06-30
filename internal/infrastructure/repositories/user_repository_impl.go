package repositories

import (
	"context"

	"time"

	"github.com/kaa-dan/clean-architecture-go/internal/domain/entities"
	"github.com/kaa-dan/clean-architecture-go/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) *UserRepository {
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

	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrUserNotFound
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil

}

func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	for cursor.Next(ctx) {
		var user entities.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, cursor.Err()
}

func (r *UserRepository) Update(ctx context.Context, id primitive.ObjectID, user *entities.User) error {
	user.UpdatedAt = time.Now()

	update := bson.M{"$set": user}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
