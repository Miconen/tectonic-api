package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Wom struct {
	Id int `json:"id"`
}

var endpoint = "https://api.wiseoldman.net/v2/players/search?username="

func GetWomId(rsn string) (string, error) {
	response, err := http.Get(endpoint + rsn)
	if err != nil {
		return "sadge", err
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

	fmt.Println(result)
	fmt.Println(result[0])

	return strconv.Itoa(result[0].Id), nil
}
