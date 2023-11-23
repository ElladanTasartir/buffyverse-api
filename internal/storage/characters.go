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

func (r *CharactersRepository) FindCharacters(ctx context.Context, params PageParams) (*PagedData[entity.Character], error) {
	orderStage := bson.D{{"$sort", bson.D{{"name", 1}}}}
	skipStage := bson.D{{"$skip", params.Page * params.PageSize}}
	limitStage := bson.D{{"$limit", params.PageSize}}
	countStage := bson.D{{"$count", "count"}}
	facet := bson.D{{"$facet", bson.D{{"results", bson.A{orderStage, skipStage, limitStage}}, {"count", bson.A{countStage}}}}}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{facet})
	if err != nil {
		return nil, fmt.Errorf("failed to find characters. err = %v", err)
	}

	var character PagedData[entity.Character]
	for cursor.Next(ctx) {
		if err := cursor.Decode(&character); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error decoding characters. err = %v", err)
		}
	}

	return &character, nil
}

func (r *CharactersRepository) CreateCharacters(ctx context.Context, characters []entity.Character) error {
	var errors []error

	for _, character := range characters {
		err := r.FindCharacterByName(ctx, &character)
		if err != nil {
			errors = append(errors, err)
		}

		fmt.Println(character.Name, character.ID)

		if character.ID == primitive.NilObjectID {
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

func (r *CharactersRepository) FindCharacterByName(ctx context.Context, character *entity.Character) error {
	filter := bson.M{
		"name": character.Name,
	}

	var foundCharacter entity.Character
	err := r.collection.FindOne(ctx, filter).Decode(&foundCharacter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		return fmt.Errorf("failed to find character. err = %v", err)
	}

	character.ID = foundCharacter.ID

	return nil
}
