package impls

import (
	"context"
	"fmt"

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

func (u *UserRepositoryImpl) CreateWalletHistory(ctx context.Context, walletHistory *model.WalletHistory) (*model.WalletHistory, error) {
	//  create wallet history
	collection := u.db.Collection("wallet_histories")
	newWalletHistory, err := collection.InsertOne(ctx, walletHistory)
	if err != nil {
		return nil, err
	}
	walletHistory.Id = newWalletHistory.InsertedID.(primitive.ObjectID).Hex()
	return walletHistory, nil
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

func (u *UserRepositoryImpl) CreditUser(ctx context.Context, amount float64, userEmail string, walletHistory *model.WalletHistory) error {
	//  update user credit
	collection := u.db.Collection("users")
	filter := bson.M{"email": userEmail}
	update := bson.M{"$inc": bson.M{"wallet_balance": amount}}
	_, err := collection.UpdateOne(ctx, filter, update)
	// create wallet history record
	if err == nil {
		_, err = u.CreateWalletHistory(ctx, walletHistory)
	}
	return err
}

func (u *UserRepositoryImpl) DebitUser(ctx context.Context, amount float64, userEmail string, walletHistory *model.WalletHistory) error {
	//  update user credit
	collection := u.db.Collection("users")
	filter := bson.M{"email": userEmail}
	update := bson.M{"$inc": bson.M{"wallet_balance": -amount}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err == nil {
		_, err = u.CreateWalletHistory(ctx, walletHistory)
	}
	return err
}

func (u *UserRepositoryImpl) UpdateUserBankDetails(ctx context.Context, userEmail string, accountDetails *model.AccountDetails) error {
	fmt.Println("Updating bank details: ", userEmail)
	fmt.Println(accountDetails)
	collection := u.db.Collection("users")
	filter := bson.M{"email": userEmail}
	update := bson.M{"$set": bson.M{"account_details": accountDetails}}
	updatedResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("Error updating bank details:", err)
		return err
	}
	if updatedResult.MatchedCount == 0 {
		fmt.Println("No user found with the provided ID.")
		return fmt.Errorf("no user found with the provided ID")
	}
	return nil
}

func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, updateMap map[string]interface{}, email string) (*model.User, error) {
	collection := u.db.Collection("users")
	filter := bson.M{"email": email}
	update := bson.M{"$set": updateMap}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return u.GetUserByEmail(ctx, email)
}

func (u *UserRepositoryImpl) GetUserByVirtualAccount(ctx context.Context, virtualAccount string) (*model.User, error) {
	collection := u.db.Collection("users")
	var user model.User
	filter := bson.M{"account_details.account_no": virtualAccount}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepositoryImpl) CreateNewFundsTracking(ctx context.Context, fundsTracking *model.FundsTracking) (*model.FundsTracking, error) {
	collection := u.db.Collection("funds_trackings")
	newFundsTracking, err := collection.InsertOne(ctx, fundsTracking)
	if err != nil {
		return nil, err
	}
	fundsTracking.Id = newFundsTracking.InsertedID.(primitive.ObjectID).Hex()
	return fundsTracking, nil
}
