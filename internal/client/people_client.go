package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreatePeople() {
	fmt.Println("Hi Ash")
	client := &http.Client{}
	endpoint := "https://api.planningcenteronline.com/people/v2/people/"
	secret_token := "a79fa5cdeef8708cfd85cc6b86f06c2540be2ae4a5f148de1782d40f9b2ec356"
	app_id := "86f5e8cd60ae4e4b8c8e48342bd86f261e270f51adeff62b090a1f69aff41302"

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
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	fmt.Println(response.Status)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	fmt.Println(string(body))
}
func main() {
	fmt.Println("Hi Ash")
	client := &http.Client{}
	endpoint := "https://api.planningcenteronline.com/people/v2/people/"
	secret_token := "a79fa5cdeef8708cfd85cc6b86f06c2540be2ae4a5f148de1782d40f9b2ec356"
	app_id := "86f5e8cd60ae4e4b8c8e48342bd86f261e270f51adeff62b090a1f69aff41302"

	goBody := Root{
		Data: Person{
			Type: "Person",
			Attributes: Attributes{
				FirstName:         "Manual",
				LastName:          "Testerdy",
				SiteAdministrator: false,
				Gender:            "male",
			},
		},
	} // Create a custom structure
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
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	fmt.Println(response.Status)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	fmt.Println(string(body))
}
