package impls

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/repository"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/skip2/go-qrcode"
)

type PartyServiceImpl struct {
	// depends on  repo
	repo         repository.PartyRepository
	azureService utils.AzureService
}

func NewPartyServiceImpl(repo repository.PartyRepository, azureService utils.AzureService) *PartyServiceImpl {
	return &PartyServiceImpl{
		repo:         repo,
		azureService: azureService,
	}
}

func (ps *PartyServiceImpl) CreateParty(createPartyDto dto.CreatePartyDTO) (map[string]interface{}, error) {
	party := model.NewParty(createPartyDto.Name, createPartyDto.Tag, createPartyDto.HostEmail)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	createdParty, err := ps.repo.CreateParty(ctx, party)
	if err != nil {
		return nil, err
	}
	inviteCode := createdParty.InviteCode
	// generate qr code with the invite code
	fileName := inviteCode + ".png"
	var qrBuffer bytes.Buffer
	nErr := qrcode.WriteFile(inviteCode, qrcode.Medium, 256, fileName)
	if nErr != nil {
		log.Fatalf("Failed to generate QR code: %v", nErr)
	}
	blobUrl, bErr := ps.azureService.UploadFileToAzureBlob(qrBuffer.Bytes(), fileName, "spraycashnew")
	respData := map[string]interface{}{
		"inviteCode": inviteCode,
	}
	if bErr != nil {
		log.Println("Error uploading qr code to azure: ", bErr)
		return respData, nil
	}
	defer func() {
		err = os.Remove(fileName)
	}()
	// delete the local file
	if err != nil {
		log.Printf("Failed to delete local file: %v", err)
	} else {
		fmt.Printf("Local file '%s' deleted successfully.\n", fileName)
	}
	respData["qrCode"] = blobUrl

	return respData, nil

}

func (ps *PartyServiceImpl) GetParty(inviteCode string) (*model.Party, error) {
	// get the party by the invite code
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	party, err := ps.repo.GetPartyByInviteCode(ctx, inviteCode)
	if err != nil {
		return nil, err
	}
	return party, nil
}

func (ps *PartyServiceImpl) JoinParty(inviteCode string, partyGuest *model.PartyGuest) (*model.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	party, err := ps.repo.JoinPartyPersist(ctx, inviteCode, partyGuest)
	if err != nil {
		return nil, err
	}
	return party, nil
}

func (ps *PartyServiceImpl) LeaveParty(inviteCode string, guestId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := ps.repo.RemoveGuest(ctx, inviteCode, guestId)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PartyServiceImpl) GetAllPartyGuests(inviteCode string) ([]*model.PartyGuest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	guests, err := ps.repo.GetAllPartyGuests(ctx, inviteCode)
	if err != nil {
		return nil, err
	}
	return guests, nil
}
