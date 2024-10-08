package mongodistrubutedlock

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
}

func SetDbConfig(client *mongo.Client, dbName, collectionName string) {
	_client = client
	_dbName = dbName
	_collectionName = collectionName
}

func New(lockKey string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		client:     _client,
		collection: _client.Database(_dbName).Collection(_collectionName),
		lockKey:    lockKey,
		expiration: expiration,
	}
}

func (d *DistributedLock) Acquire(ctx context.Context) error {
	now := time.Now()
	expirationTime := now.Add(d.expiration)

	_, err := d.collection.InsertOne(ctx, bson.M{
		"_id":        d.lockKey,
		"expires_at": expirationTime,
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
