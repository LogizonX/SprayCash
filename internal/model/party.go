package model

import (
	"fmt"
	"strings"
	"sync"

	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/gorilla/websocket"
)

type Party struct {
	Id         string `bson:"_id,omitempty" json:"id"`
	Name       string `bson:"name" json:"name"`
	Tag        string `bson:"tag" json:"tag"`
	HostEmail  string `bson:"hostEmail" json:"hostEmail"`
	InviteCode string `bson:"inviteCode" json:"inviteCode"`
	mu         sync.RWMutex
	Guests     map[string]*PartyGuest `bson:"guests" json:"guests"`
}

// TODO: Guest total spent should be updated when a user makes a payment
// TODO: Guest total received should updated

type PartyConnPool struct {
	mu     sync.RWMutex
	Guests map[string]*PartyGuest
	Party  *Party
}

// party broadcast message
// from this, I will loop through the clients and write the message
func (p *PartyConnPool) BroadcastMessage(msg *Message) {
	fmt.Println("Broadcasting message to party guests")
	fmt.Println("pool guest: %w", p.Guests)
	for _, guest := range p.Guests {
		guest.conn.WriteJSON(msg)
		fmt.Println("Broadcasting message to ", guest.Username)
	}
}

// Implementing the singleton pattern, to ensure that there is only one instance of the party connection pool
var (
	instance *PartyConnPool
	once     sync.Once
)

func NewPartyConnPool(party *Party) *PartyConnPool {
	once.Do(func() {
		instance = &PartyConnPool{
			Guests: make(map[string]*PartyGuest),
			Party:  party,
		}
	})
	return instance
}

// broadcast ranking
func (p *PartyConnPool) BroadcastRanking() {
	for _, guest := range p.Guests {
		guest.conn.WriteJSON(p.GetRanking())
	}
}

func (p *PartyConnPool) GetRanking() []*GuestsData {
	var ranking []*GuestsData
	for _, guest := range p.Guests {
		ranking = append(ranking, NewGuestData(guest.PartyId, guest.UserId, guest.Username, int64(guest.TotalSpent)))
	}
	return ranking
}

// leave party
func (p *PartyConnPool) LeaveParty(userId string) {
	guest := p.Guests[userId]
	delete(p.Guests, userId)
	go p.BroadcastMessage(NewMessage(p.Party.Id, fmt.Sprintf("%s has left the party", guest.Username), guest.Username, guest.UserId))
}

// join party
func (p *PartyConnPool) JoinParty(guest *PartyGuest) {
	fmt.Println("Adding user: ", guest.UserId)
	p.Guests[guest.UserId] = guest
	// broadcast a new user joining the party
	go p.BroadcastMessage(NewMessage(p.Party.Id, fmt.Sprintf("%s has joined the party", guest.Username), guest.Username, guest.UserId))
}

// read message from client

// broadcast message to client

// PartyGuest, the websocket client representation
type PartyGuest struct {
	PartyId        string `bson:"partyId" json:"partyId"`
	UserId         string `bson:"userId" json:"userId"`
	Email          string `bson:"email" json:"email"`
	CanReceiveFund bool   `bson:"canReceiveFund" json:"canReceiveFund"`
	Username       string `bson:"username" json:"username"`
	TotalSpent     int64  `bson:"totalSpent" json:"totalSpent"`
	conn           *websocket.Conn
}

type GuestsData struct {
	PartyId    string `bson:"partyId" json:"partyId"`
	UserId     string `bson:"userId" json:"userId"`
	Username   string `bson:"username" json:"username"`
	TotalSpent int64  `bson:"totalSpent" json:"totalSpent"`
}

func NewGuestData(partyId string, userId string, username string, totalSpent int64) *GuestsData {
	return &GuestsData{
		PartyId:    partyId,
		UserId:     userId,
		Username:   username,
		TotalSpent: int64(totalSpent),
	}
}

type MessageData struct {
	ReceiverEmail string `json:"receiverEmail"`
	ReceiverName  string `json:"receiverName"`
	SenderEmail   string `json:"senderEmail"`
	SenderName    string `json:"senderName"`
	Amount        int64  `json:"amount"`
}

type Message struct {
	PartyId  string `json:"partyId"`
	Message  string `json:"message"`
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

func NewParty(name string, tag string, hostEmail string) *Party {
	inviteCode := strings.ReplaceAll(name, " ", "-") + "-" + utils.GenerateInviteCode()
	return &Party{
		Name:       name,
		Tag:        tag,
		HostEmail:  hostEmail,
		InviteCode: inviteCode,
		Guests:     make(map[string]*PartyGuest),
		mu:         sync.RWMutex{},
	}
}

func NewMessage(partyId string, message string, username string, userId string) *Message {
	return &Message{
		PartyId:  partyId,
		Message:  message,
		Username: username,
		UserId:   userId,
	}
}

func NewPartyGuest(partyId string, email string, conn *websocket.Conn, userId string) *PartyGuest {
	return &PartyGuest{
		PartyId:        partyId,
		Email:          email,
		CanReceiveFund: false,
		TotalSpent:     0,
		conn:           conn,
		UserId:         userId,
		Username:       strings.Split(email, "@")[0],
	}
}
