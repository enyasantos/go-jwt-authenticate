package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"password,omitempty" bson:"password"`
	Email     string        `json:"email" bson:"email"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

func (u *User) Validate() error {
	errorMessage := []string{}

	if len(u.Username) < 6 {
		errorMessage = append(errorMessage, "Username must be at least 6 characters long")
	}

	if len(u.Email) < 6 {
		errorMessage = append(errorMessage, "Email must be at least 6 characters long")
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		errorMessage = append(errorMessage, "Invalid email format")
	}

	if len(u.Password) < 8 {
		errorMessage = append(errorMessage, "Password must be at least 8 characters long")
	}

	if len(errorMessage) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMessage, ","))
	}

	return nil
}

func NewUser(username, email, password string) (*User, error) {
	user := &User{
		ID:        bson.NewObjectID(),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.Hex(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
