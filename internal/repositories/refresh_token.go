package repositories

import (
	"authentication-jwt/internal/database"
	"authentication-jwt/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RefreshTokenRepositoryInterface interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	Delete(ctx context.Context, token string) error
}

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *database.Database) *RefreshTokenRepository {
	collection := db.Client.Collection("refresh_tokens")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "token", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(fmt.Sprintf("Failed to create index on refresh_tokens collection: %v", err))
	}

	return &RefreshTokenRepository{
		collection: collection,
	}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	filter := bson.M{"user_id": token.UserID}
	update := bson.M{
		"$set": bson.M{
			"token":      token.Token,
			"created_at": token.CreatedAt,
			"expires_at": token.ExpiresAt,
			"updated_at": token.UpdatedAt,
		},
	}
	opts := options.UpdateOne().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&refreshToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No refresh token found for the user
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return err
	}

	return nil
}
