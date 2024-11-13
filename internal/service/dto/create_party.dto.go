package dto

type CreatePartyDTO struct {
	Name      string `json:"name"`
	HostEmail string `json:"host_email"`
	Tag       string `json:"tag"`
}
