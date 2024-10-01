package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sen1or/micromovie/services/movie/internal/gateway"
	model "sen1or/micromovie/services/rating/pkg"
)

type Gateway struct {
	addr string
}

func NewGateway(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
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
	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", nil)
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
