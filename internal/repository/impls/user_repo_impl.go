package impls

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	//  depends on mongodb database
	db *mongo.Database
}

func NewUserRepositoryImpl(db *mongo.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	collection := u.db.Collection("users")
	newUser, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = newUser.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}
