package impls

import (
	"context"
	"time"

	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/repository"
	"github.com/LoginX/SprayDash/internal/service/dto"
)

type PartyServiceImpl struct {
	// depends on  repo
	repo repository.PartyRepository
}

func NewPartyServiceImpl(repo repository.PartyRepository) *PartyServiceImpl {
	return &PartyServiceImpl{
		repo: repo,
	}
}

func (ps *PartyServiceImpl) CreateParty(createPartyDto dto.CreatePartyDTO) (*model.Party, error) {
	party := model.NewParty(createPartyDto.Name, createPartyDto.Tag, createPartyDto.HostEmail)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	createdParty, err := ps.repo.CreateParty(ctx, party)
	if err != nil {
		return nil, err
	}
	return createdParty, nil

}

func (ps *PartyServiceImpl) JoinParty(inviteCode string) (*model.Party, error) {
	// get the party by the invite code
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	party, err := ps.repo.GetPartyByInviteCode(ctx, inviteCode)
	if err != nil {
		return nil, err
	}
	return party, nil
}
