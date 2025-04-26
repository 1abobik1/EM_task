package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/1abobik1/EM_task/internal/dto"
	"github.com/sirupsen/logrus"
)

type NationalizeClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewNationalizeClient(baseURL string, httpClient *http.Client) *NationalizeClient {
	return &NationalizeClient{baseURL: baseURL, httpClient: httpClient}
}

func (c *NationalizeClient) GetNationality(ctx context.Context, name string) (string, error) {
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

	var natResp dto.NatResp

	if err := json.NewDecoder(resp.Body).Decode(&natResp); err != nil {
		logrus.Error(err)
		return "", nil
	}

	if len(natResp.Country) == 0 {
		return "", errors.New("no nationality data")
	}

	max := natResp.Country[0]
	for _, c := range natResp.Country {
		if c.Probability > max.Probability {
			max = c
		}
	}

	return max.CountryID, nil
}
