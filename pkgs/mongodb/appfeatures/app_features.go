package appfeatures

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppFeatures struct {
	appName  string
	features map[string]AppFeature
	locker   sync.RWMutex

	_client         *mongo.Client
	_dbName         string
	_collectionName string
}

type AppFeature interface {
	Upgrade() error
}

func New(client *mongo.Client, dbName, collectionName, appName string) *AppFeatures {
	return &AppFeatures{
		_client:         client,
		_dbName:         dbName,
		_collectionName: collectionName,
		appName:         appName,

		features: map[string]AppFeature{},
	}
}

func (af *AppFeatures) EnableFeaturesAndUpgrade(ctx context.Context, features map[string]AppFeature) error {
	af.locker.Lock()
	defer af.locker.Unlock()

	for featureName, feature := range features {
		if _, ok := af.features[featureName]; ok {
			continue
		}

		if feature != nil {
			if err := feature.Upgrade(); err != nil {
				return fmt.Errorf("failed to upgrade feautre %s, err: %w", featureName, err)
			}
		}

		af.features[featureName] = feature
		if err := af.enableDbFeature(ctx, af.appName, featureName); err != nil {
			return err
		}
	}

	return nil
}

func (af *AppFeatures) Init(ctx context.Context) error {
	if err := af.createIndexes(ctx); err != nil {
		return err
	}
	if err := af.createRecordIfNotExisted(ctx); err != nil {
		return err
	}

	return nil
}
