package service

import (
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/service/dto"
)

type PartyService interface {
	CreateParty(createPartyDto dto.CreatePartyDTO) (*model.Party, error)
}
