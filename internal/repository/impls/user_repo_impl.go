package impls

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
	"go.mongodb.org/mongo-driver/bson"
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

func (u *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	collection := u.db.Collection("users")
	var user model.User
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepositoryImpl) CreditUser(ctx context.Context, amount float64, userId string) error {
	//  update user credit
	collection := u.db.Collection("users")
	filter := bson.M{"_id": userId}
	update := bson.M{"$inc": bson.M{"wallet_balance": amount}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (u *UserRepositoryImpl) DebitUser(ctx context.Context, amount float64, userId string) error {
	//  update user credit
	collection := u.db.Collection("users")
	filter := bson.M{"_id": userId}
	update := bson.M{"$inc": bson.M{"wallet_balance": -amount}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
