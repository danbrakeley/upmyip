package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PublicInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	IP          string  `json:"query"`
}

// this is an order of magnitude larger than the responses I got in my local testing,
// so should be a good upper bound for any possible valid response.
const maxGetIPResponseSize = 4096

func RequestPublicInfo(ctx context.Context) (*PublicInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://ip-api.com/json/", nil)
	if err != nil {
		return nil, fmt.Errorf("new request ip-api.com: %w", err)
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, maxGetIPResponseSize))
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if len(body) >= maxGetIPResponseSize {
		return nil, fmt.Errorf("response from ip-api.com was longer than expected")
	}

	var resp PublicInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling body: %w", err)
	}

	return &resp, nil
}
