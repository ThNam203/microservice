package http

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sen1or/micromovie/movie/internal/gateway"
	"sen1or/micromovie/pkg/discovery"
	model "sen1or/micromovie/rating/pkg"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	url, err := getURLFromRegistry(ctx, g.registry)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	req.URL.RawQuery = values.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx status code: %s", err)
	}

	var sum float64
	if err := json.NewDecoder(res.Body).Decode(&sum); err != nil {
		return 0, err
	}

	return sum, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	url, err := getURLFromRegistry(ctx, g.registry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}

func getURLFromRegistry(ctx context.Context, registry discovery.Registry) (string, error) {
	addrs, err := registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return "", nil
	}

	url := addrs[rand.Intn(len(addrs))] + "/rating"
	if err != nil {
		return "", err
	}

	return url, nil
}
