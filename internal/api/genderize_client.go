package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1abobik1/EM_task/internal/dto"
	"github.com/sirupsen/logrus"
)

type GenderizeClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewGenderizeClient(baseURL string, httpClient *http.Client) *GenderizeClient {
	return &GenderizeClient{baseURL: baseURL, httpClient: httpClient}
}

func (c *GenderizeClient) GetGender(ctx context.Context, name string) (string, error) {
	queryURL := fmt.Sprintf("%s/?name=%s", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	var genderResp dto.GenderResp
	
	if err := json.NewDecoder(resp.Body).Decode(&genderResp); err != nil {
		logrus.Error(err)
		return "", err
	}

	return genderResp.Gender, nil
}
