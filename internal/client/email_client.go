package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

)

func GetEmail(client *PC_Client, app_id, secret_token, id string) EmailRoot {
	//Fetch the data
	endpoint := HostURL + "people/v2/emails/" + id
	request, err := http.NewRequest("GET", endpoint, nil)

	request.SetBasicAuth(app_id, secret_token)

	if err != nil {
		fmt.Println("Error:", err)
	}
	response, err := client.Client.Do(request)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var jsonBody EmailRoot
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	return jsonBody

}

// '{"data": {"type": "Email", "attributes": {"address": "tester@hcubedcoder.com", "location": "home", "primary": true}}}'
func CreateEmail(client *PC_Client, app_id, secret_token, peopleID string, responseData *EmailRoot) []byte {
	endpoint := HostURL + "people/v2/people/" + peopleID + "/emails"

	// Convert struct to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	// Create a request with the JSON data
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set the content type to application/json
	request.Header.Set("Content-Type", "application/json")

	// Make the request
	request.SetBasicAuth(app_id, secret_token)
	response, err := client.Client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	return body
}

func DeleteEmail(client *PC_Client, app_id, secret_token, id string) {
	endpoint := HostURL + "people/v2/emails/" + id

	// Create a request with the JSON data
	request, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	request.SetBasicAuth(app_id, secret_token)
	response, err := client.Client.Do(request)

	if err != nil {
		fmt.Println("Error sending request: ", err)
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	fmt.Println(string(body))

}

func UpdateEmail(client *PC_Client, app_id, secret_token, id string, responseData *EmailRoot) []byte {
	endpoint := HostURL + "people/v2/emails/" + id

	// Convert struct to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	// Create a request with the JSON data
	request, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set the content type to application/json
	request.Header.Set("Content-Type", "application/json")

	// Make the request
	request.SetBasicAuth(app_id, secret_token)
	response, err := client.Client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	return body

}
