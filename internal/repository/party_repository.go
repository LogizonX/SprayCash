package repository

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
)

type PartyRepository interface {
	CreateParty(ctx context.Context, party *model.Party) (*model.Party, error)
}
