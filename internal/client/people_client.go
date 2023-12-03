package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreatePeople(client PC_Client, app_id, secret_token string) {
	fmt.Println("Hi Ash")
  endpoint := HostURL + "/people/v2/people/"

	var goBody Root
	// Create a custom structure
	customStruct := struct {
		Data Person `json:"data"`
	}{goBody.Data}

	// Convert struct to JSON
	jsonData, err := json.MarshalIndent(customStruct, "", " ")
	fmt.Println("*******************************************")
	fmt.Printf("%s\n", jsonData)
	fmt.Println("*******************************************")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create a request with the JSON data
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type to application/json
	request.Header.Set("Content-Type", "application/json")

	// Make the request
	request.SetBasicAuth(app_id, secret_token)
	response, err := client.Client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	fmt.Println(response.Status)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	fmt.Println(string(body))
}

