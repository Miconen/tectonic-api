package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tectonic-api/config"
	"tectonic-api/logging"
	"tectonic-api/models"
)

type Wom struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

type WomClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewWomClient(cfg *config.Config) *WomClient {
	return &WomClient{
		baseURL: cfg.WOM.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.WOM.Timeout,
		},
	}
}

func (c *WomClient) GetWom(rsn string) (Wom, error) {
	url := c.baseURL + "/players/" + rsn
	return handleResponse[Wom](url, c)
}

func (c *WomClient) GetCompetition(id int) (models.WomCompetition, error) {
	url := c.baseURL + "/competitions/" + strconv.Itoa(id)
	return handleResponse[models.WomCompetition](url, c)
}

func handleResponse[T any](url string, c *WomClient) (T, error) {
	var result T

	ctx, cancel := context.WithTimeout(context.Background(), c.httpClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logging.Get().Error("failed to create wom api request", "url", url, "error", err)
		return result, err
	}

	// Send request and handle timeout errors
	response, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logging.Get().Error("wom request timed out", "url", url, "timeout", c.httpClient.Timeout, "error", err)
		} else {
			logging.Get().Error("wom request failed", "url", url, "error", err)
		}
		return result, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logging.Get().Error("unexpected wom api status code", "status", response.StatusCode)
		return result, errors.New("Unexpected status code:" + strconv.Itoa(response.StatusCode))
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		logging.Get().Error("failed to decode wom response", "error", err)
		fmt.Println("Error decoding JSON:", err)
		return result, err
	}

	return result, nil
}
