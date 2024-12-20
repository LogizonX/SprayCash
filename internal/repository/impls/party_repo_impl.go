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

func (p *PartyRepoImpl) GetPartyByInviteCode(ctx context.Context, inviteCode string) (*model.Party, error) {
	collection := p.db.Collection("party")
	var party model.Party
	err := collection.FindOne(ctx, primitive.M{"inviteCode": inviteCode}).Decode(&party)
	if err != nil {
		return nil, err
	}
	return &party, nil
}

func (p *PartyRepoImpl) JoinPartyPersist(ctx context.Context, inviteCode string, partyGuest *model.PartyGuest) (*model.Party, error) {
	collection := p.db.Collection("party")
	_, err := collection.UpdateOne(ctx, primitive.M{"inviteCode": inviteCode}, primitive.M{"$set": primitive.M{"guests." + partyGuest.UserId: partyGuest}})
	if err != nil {
		return nil, err
	}
	var updatedParty model.Party
	err = collection.FindOne(ctx, primitive.M{"inviteCode": inviteCode}).Decode(&updatedParty)
	if err != nil {
		return nil, err
	}

	return &updatedParty, nil
}

func (p *PartyRepoImpl) RemoveGuest(ctx context.Context, inviteCode string, guestId string) error {
	collection := p.db.Collection("party")
	_, err := collection.UpdateOne(ctx, primitive.M{"inviteCode": inviteCode}, primitive.M{"$unset": primitive.M{"guests." + guestId: ""}})
	if err != nil {
		return err
	}
	return nil
}

func (p *PartyRepoImpl) GetAllPartyGuests(ctx context.Context, inviteCode string) ([]*model.PartyGuest, error) {
	collection := p.db.Collection("party")
	var party model.Party
	err := collection.FindOne(ctx, primitive.M{"inviteCode": inviteCode}).Decode(&party)
	if err != nil {
		return nil, err
	}
	var guests []*model.PartyGuest
	for _, guest := range party.Guests {
		guests = append(guests, guest)
	}
	return guests, nil
}
