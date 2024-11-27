package service

import (
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/service/dto"
)

type PartyService interface {
	CreateParty(createPartyDto dto.CreatePartyDTO) (map[string]interface{}, error)
	GetParty(inviteCode string) (*model.Party, error)
	JoinParty(inviteCode string, partyGuest *model.PartyGuest) (*model.Party, error)
	LeaveParty(inviteCode string, guestId string) error
	GetAllPartyGuests(inviteCode string) ([]*model.PartyGuest, error)
}
