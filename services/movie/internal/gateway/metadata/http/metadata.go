package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	model "sen1or/micromovie/services/metadata/pkg"
	"sen1or/micromovie/services/movie/internal/gateway"
)

type Gateway struct {
	addr string
}

func NewGateway(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	request, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
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
