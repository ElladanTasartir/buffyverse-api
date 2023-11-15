package storage

import (
	"context"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
)

type CharactersStorage interface {
	CreateCharacters(ctx context.Context, characters []entity.Character) error
	FindCharacterByName(ctx context.Context, character entity.Character) (*entity.Character, error)
}
