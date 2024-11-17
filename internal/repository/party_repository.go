package repository

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
)

type PartyRepository interface {
	CreateParty(ctx context.Context, party *model.Party) (*model.Party, error)
	GetPartyByInviteCode(ctx context.Context, inviteCode string) (*model.Party, error)
	JoinPartyPersist(ctx context.Context, inviteCode string, partyGuest *model.PartyGuest) (*model.Party, error)
	RemoveGuest(ctx context.Context, inviteCode string, guestId string) error
	GetAllPartyGuests(ctx context.Context, inviteCode string) ([]*model.PartyGuest, error)
}
