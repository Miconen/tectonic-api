package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Wom struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

var endpoint = "https://api.wiseoldman.net/v2/players/search?username="

func GetWomId(rsn string) (string, error) {
	response, err := http.Get(endpoint + rsn)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", response.StatusCode)
		return "", errors.New("Unexpected status code:" + strconv.Itoa(response.StatusCode))
	}

	var result []Wom
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	if len(result) == 0 {
		fmt.Println("RSN not found")
		return "", errors.New("RSN not found")
	}

	id := strconv.Itoa(result[0].Id)
	name := result[0].DisplayName
	if name != rsn {
		fmt.Printf("Provided RSN (%s) doesn't match fetched RSN (%s)\n", rsn, name)
		return "", errors.New("Provided RSN doesn't match fetched RSN")
	}

	return id, nil
}
