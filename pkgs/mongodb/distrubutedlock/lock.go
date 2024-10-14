package distrubutedlock

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrLockAcquired    = errors.New("lock already acquired")
	ErrLockNotAcquired = errors.New("lock not acquired")
)

var (
	_client         *mongo.Client
	_dbName         string
	_collectionName string
)

type DistributedLock struct {
	client     *mongo.Client
	collection *mongo.Collection
	lockKey    string
	expiration time.Duration
	labels     map[string]string
}

func Init(ctx context.Context, client *mongo.Client, dbName, collectionName string) error {
	_client = client
	_dbName = dbName
	_collectionName = collectionName

	return createIndexes(ctx)
}

func createIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys:    bson.M{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := _client.Database(_dbName).Collection(_collectionName).Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

func New(lockKey string, expiration time.Duration, labels map[string]string) *DistributedLock {
	return &DistributedLock{
		client:     _client,
		collection: _client.Database(_dbName).Collection(_collectionName),
		lockKey:    lockKey,
		expiration: expiration,
		labels:     labels,
	}
}

func (d *DistributedLock) Acquire(ctx context.Context) error {
	now := time.Now()
	expirationTime := now.Add(d.expiration)

	_, err := d.collection.InsertOne(ctx, bson.M{
		"_id":        d.lockKey,
		"expires_at": expirationTime,
		"labels":     d.labels,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrLockAcquired
		}
		return err
	}

	return nil
}

func (d *DistributedLock) Release(ctx context.Context) error {
	_, err := d.collection.DeleteOne(ctx, bson.M{
		"_id": d.lockKey,
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *DistributedLock) IsLocked(ctx context.Context) (bool, error) {
	count, err := d.collection.CountDocuments(ctx, bson.M{
		"_id": d.lockKey,
	})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d *DistributedLock) GetLabels(ctx context.Context) (map[string]string, error) {
	var result struct {
		Labels map[string]string `bson:"labels"`
	}
	err := d.collection.FindOne(ctx, bson.M{
		"_id": d.lockKey,
	}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result.Labels, nil
}
