package appfeatures

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type appFeaturesEntity struct {
	ID       *primitive.ObjectID `bson:"_id,omitempty"`
	App      string              `bson:"app"`
	Features []string            `bson:"features"`
}

func (af *AppFeatures) createIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "app", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	collection := af._client.Database(af._dbName).Collection(af._collectionName)
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (af *AppFeatures) createRecordIfNotExisted(ctx context.Context) error {
	entity := appFeaturesEntity{
		App:      af.appName,
		Features: []string{},
	}

	filter := bson.D{{Key: "app", Value: af.appName}}

	collection := af._client.Database(af._dbName).Collection(af._collectionName)

	var existingEntity appFeaturesEntity
	err := collection.FindOne(ctx, filter).Decode(&existingEntity)
	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, entity)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (af *AppFeatures) enableDbFeature(ctx context.Context, appName, featureName string) error {
	filter := bson.D{{Key: "app", Value: appName}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "features", Value: featureName},
		}},
	}

	collection := af._client.Database(af._dbName).Collection(af._collectionName)
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
