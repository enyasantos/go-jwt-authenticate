package repositories

import (
	"authentication-jwt/internal/database"
	"authentication-jwt/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindById(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		collection: db.Client.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	log.Default().Printf("Finding user with ID: %s", objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No user found
		}
		return nil, err
	}

	log.Default().Printf("User found: %+v", user)

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No user found
		}
		return nil, err // Other errors
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	if user.ID.IsZero() {
		return errors.New("user ID is required for update")
	}

	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	if err != nil {
		return err
	}

	return nil
}
