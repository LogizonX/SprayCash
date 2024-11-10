package model

import "github.com/LoginX/SprayDash/internal/utils"

type Party struct {
	Id         string `bson:"_id,omitempty" json:"id"`
	Name       string `bson:"name" json:"name"`
	Tag        string `bson:"tag" json:"tag"`
	HostEmail  string `bson:"hostEmail" json:"hostEmail"`
	InviteCode string `bson:"inviteCode" json:"inviteCode"`
}

type PartyGuest struct {
	PartyId string `bson:"partyId" json:"partyId"`
	Email   string `bson:"email" json:"email"`
}

func NewParty(name string, tag string, hostEmail string) *Party {
	inviteCode := name + "-" + utils.GenerateInviteCode()
	return &Party{
		Name:       name,
		Tag:        tag,
		HostEmail:  hostEmail,
		InviteCode: inviteCode,
	}
}

func NewPartyGuest(partyId string, email string) *PartyGuest {
	return &PartyGuest{
		PartyId: partyId,
		Email:   email,
	}
}
