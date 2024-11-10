package impls

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyRepoImpl struct {
	// depends on
	db *mongo.Database
}

func NewPartyRepoImpl(db *mongo.Database) *PartyRepoImpl {
	return &PartyRepoImpl{
		db: db,
	}
}

func (p *PartyRepoImpl) CreateParty(ctx context.Context, party *model.Party) (*model.Party, error) {
	collection := p.db.Collection("party")
	result, err := collection.InsertOne(ctx, party)
	if err != nil {
		return nil, err
	}
	// return the invite code
	party.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return party, nil

}
