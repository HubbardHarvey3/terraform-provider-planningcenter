package client

import (
	"fmt"
	"net/http"
)

const HostURL = "https://api.planningcenteronline.com/"

type PC_Client struct {
	Client   *http.Client
	Token    string
	AppID    string
	Endpoint string
}

func NewPCClient(id, token, endpoint string) *PC_Client {
	fmt.Println("Returning a new PCClient")
	return &PC_Client{
		Client:   &http.Client{},
		AppID:    id,
		Token:    token,
		Endpoint: endpoint,
	}
}
