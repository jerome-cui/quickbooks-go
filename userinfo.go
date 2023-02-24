package quickbooks

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Address struct {
	StreetAddress string `json:"streetAddress"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
}

type UserInfoResponse struct {
	Sub                 string  `json:"sub"`
	Email               string  `json:"email"`
	EmailVerified       bool    `json:"emailVerified"`
	GivenName           string  `json:"givenName"`
	FamilyName          string  `json:"familyName"`
	PhoneNumber         string  `json:"phoneNumber"`
	PhoneNumberVerified bool    `json:"phoneNumberVerified"`
	Address             Address `json:"address"`
}

/*
 * Method to retrive userInfo - email, address, name, phone etc
 */
func (c *Client) GetUserInfo(accessToken string) (*UserInfoResponse, error) {
	log.Println("Inside GetUserInfo ")
	client := &http.Client{}

	request, err := http.NewRequest("GET", c.discoveryAPI.UserinfoEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//set header
	request.Header.Set("accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	userInfoResponse, err := getUserInfoResponse([]byte(body))

	log.Println("Ending GetUserInfo")
	return userInfoResponse, err
}

func getUserInfoResponse(body []byte) (*UserInfoResponse, error) {
	var s = new(UserInfoResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln("error parsing userInfoResponse:", err)
	}
	return s, err
}
