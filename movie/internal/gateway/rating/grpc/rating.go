package grpc

import (
	"context"
	"sen1or/micromovie/gen"
	"sen1or/micromovie/movie/internal/grpcutil"
	"sen1or/micromovie/pkg/discovery"
	"sen1or/micromovie/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	res, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})

	if err != nil {
		return 0, err
	}

	return float64(res.RatingValue), nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	_, err = client.PutRating(ctx, &gen.PutRatingRequest{UserId: string(rating.UserID), RecordId: string(recordID), RecordType: string(recordType), RatingValue: int32(rating.Value)})

	if err != nil {
		return err
	}

	return nil
}
