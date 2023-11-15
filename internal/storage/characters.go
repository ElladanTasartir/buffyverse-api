package storage

import (
	"context"
	"fmt"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharactersRepository struct {
	collection *mongo.Collection
}

const collection = "characters"

func NewCharactersRepository(storage *Storage) (*CharactersRepository, error) {
	colletion := storage.database.Collection(collection)

	return &CharactersRepository{
		collection: colletion,
	}, nil
}

func (r *CharactersRepository) CreateCharacters(ctx context.Context, characters []entity.Character) error {
	var errors []error

	for _, character := range characters {
		foundCharacter, err := r.FindCharacterByName(ctx, character)
		if err != nil {
			errors = append(errors, err)
		}

		if foundCharacter != nil {
			character.ID = foundCharacter.ID
		} else {
			character.ID = primitive.NewObjectID()
		}

		_, err = r.collection.InsertOne(ctx, character)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		var errorMessages string
		for i, err := range errors {
			if i == len(errors)+1 {
				errorMessages = err.Error()
				continue
			}

			errorMessages = fmt.Sprintf("%s ", err.Error())
		}

		return fmt.Errorf("there was a problem inserting scraped data. errors = %v", errorMessages)
	}

	return nil
}

func (r *CharactersRepository) FindCharacterByName(ctx context.Context, character entity.Character) (*entity.Character, error) {
	filter := bson.M{
		"name": character.Name,
	}

	var foundCharacter entity.Character
	err := r.collection.FindOne(ctx, filter).Decode(&foundCharacter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find character. err = %v", err)
	}

	return &foundCharacter, nil
}
