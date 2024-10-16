package http

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	model "sen1or/micromovie/metadata/pkg"
	"sen1or/micromovie/movie/internal/gateway"
	"sen1or/micromovie/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, nil
	}

	url := addrs[rand.Intn(len(addrs))] + "/metadata"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	values := request.URL.Query()
	values.Add("id", id)
	request.URL.RawQuery = values.Encode()

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx status code: %v", err)
	}

	var movieMetadata *model.Metadata
	if err := json.NewDecoder(res.Body).Decode(movieMetadata); err != nil {
		return nil, err
	}
	return movieMetadata, nil
}
