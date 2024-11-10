package controller

import (
	"net/http"

	"github.com/LoginX/SprayDash/internal/service"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/gin-gonic/gin"
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
