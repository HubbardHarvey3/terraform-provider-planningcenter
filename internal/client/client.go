package client

import (
	"fmt"
	"net/http"
)

const HostURL = "https://api.planningcenteronline.com/"

type PC_Client struct {
	Client   *http.Client
	id       string
	token    string
	endpoint string
}

func NewPCClient(id, token, endpoint string) *PC_Client {
  fmt.Println("Returning a new PCClient")
	return &PC_Client{
		Client:   &http.Client{},
		id:       id,
		token:    token,
		endpoint: endpoint,
	}
}
