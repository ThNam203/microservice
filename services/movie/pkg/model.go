package model

import model "sen1or/micromovie/services/metadata/pkg"

type MovieDetails struct {
	Rating   *float64       `json:"rating,omitEmpty"`
	Metadata model.Metadata `json:"metadata"`
}
