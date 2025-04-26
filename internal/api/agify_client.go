package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1abobik1/EM_task/internal/dto"
	"github.com/sirupsen/logrus"
)

type AgifyClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAgifyClient(baseURL string, httpClient *http.Client) *AgifyClient {
	return &AgifyClient{baseURL: baseURL, httpClient: httpClient}
}

func (c *AgifyClient) GetAge(ctx context.Context, name string) (int, error) {
	queryURL := fmt.Sprintf("%s/?name=%s", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Warnf("agify returned status: %v", resp.Status)
		return 0, fmt.Errorf("agify returned status: %v", resp.Status)
	}

	var ageResp dto.AgeResp

	if err := json.NewDecoder(resp.Body).Decode(&ageResp); err != nil {
		logrus.Error(err)
		return 0, err
	}

	return ageResp.Age, nil
}
