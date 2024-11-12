package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/service"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PartyController struct {
	partyService service.PartyService
}

func NewPartyController(partyService service.PartyService) PartyController {
	return PartyController{
		partyService: partyService,
	}
}

func (ps *PartyController) RegisterPartyRoutes(rg *gin.RouterGroup) {
	partyRoutes := rg.Group("/parties")
	{
		partyRoutes.POST("/create", ps.CreateParty)
	}
}

func (ps *PartyController) CreateParty(c *gin.Context) {
	var pl dto.CreatePartyDTO
	// Bind the request body to the CreatePartyDTO struct
	if err := c.ShouldBindJSON(&pl); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, nil, err.Error()))
		return
	}

	party, err := ps.partyService.CreateParty(pl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, nil, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.Response(http.StatusCreated, party, "Party created successfully"))

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ps *PartyController) JoinParty(c *gin.Context) {
	// get user from the context
	user, uErr := utils.GetUserFromContext(c)
	if uErr != nil {
		log.Println(uErr)
		c.JSON(http.StatusUnauthorized, utils.Response(http.StatusUnauthorized, nil, "Unauthorized"))
		return
	}
	// get the invite code from the query params
	inviteCode := c.Query("inviteCode")
	// get the party by the invite code
	party, err := ps.partyService.JoinParty(inviteCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, nil, err.Error()))
		return
	}
	conn, cErr := upgrader.Upgrade(c.Writer, c.Request, nil)
	if cErr != nil {
		log.Println("failed to upgrade connection: ", cErr)
		return
	}
	partyGuest := model.NewPartyGuest(party.Id, user.Email, conn, user.Id)
	party.JoinParty(partyGuest)
	// listen for message in a goroutine
	go func() {
		defer conn.Close()
		for {
			var messageData model.MessageData
			err := conn.ReadJSON(&messageData)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Unexpected close error: %v", err)
				}
				log.Println("failed to read message: ", err)
				break
			}
			// broadcast the message to all guests in the party
			message := fmt.Sprintf("%s")
			party.BroadcastMessage(message)
		}
	}()

}
