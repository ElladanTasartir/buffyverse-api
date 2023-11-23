package storage

import (
	"context"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
)

type PageParams struct {
	Page     int32
	PageSize int32
	Search   string
	Order    string
}

type PagedData[T interface{}] struct {
	Results []T `json:"results" bson:"results"`
	Count   []struct {
		Count int `json:"count" bson:"count"`
	} `json:"count"  bson:"count"`
}

type CharactersStorage interface {
	CreateCharacters(ctx context.Context, characters []entity.Character) error
	FindCharacterByName(ctx context.Context, character *entity.Character) error
	FindCharacters(ctx context.Context, params PageParams) (*PagedData[entity.Character], error)
}
