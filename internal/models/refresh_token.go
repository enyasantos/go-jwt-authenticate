package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshToken struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    bson.ObjectID `json:"user_id" bson:"user_id"`
	Token     string        `json:"token" bson:"token"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	ExpiresAt time.Time     `json:"expires_at" bson:"expires_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"expires_at"`
}

func NewRefreshToken(token string, userID bson.ObjectID, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		ID:        bson.NewObjectID(),
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		UpdatedAt: time.Now(),
	}
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}
