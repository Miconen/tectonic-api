package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tectonic-api/models"
)

type Wom struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

var endpoint = "https://api.wiseoldman.net/v2"
var players = endpoint + "/players/"
var competitions = endpoint + "/competitions/"

func GetWom(rsn string) (Wom, error) {
	var result Wom

	response, err := http.Get(players + rsn)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", response.StatusCode)
		return result, errors.New("Unexpected status code:" + strconv.Itoa(response.StatusCode))
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return result, err
	}

	return result, nil
}

func GetCompetition(id int) (models.WomCompetition, error) {
	var result models.WomCompetition

	response, err := http.Get(competitions + strconv.Itoa(id))
	if err != nil {
		return result, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", response.StatusCode)
		return result, errors.New("Unexpected status code:" + strconv.Itoa(response.StatusCode))
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return result, err
	}

	return result, nil
}
