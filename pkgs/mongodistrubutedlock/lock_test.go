package mongodistrubutedlock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func TestDistributedLock(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	defer mongoC.Terminate(ctx)

	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + port.Port())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}

	SetDbConfig(client, "testdb", "locks")

	lock := New("test_lock", 5*time.Minute)

	err = lock.Acquire(ctx)
	assert.NoError(t, err)

	isLocked, err := lock.IsLocked(ctx)
	assert.NoError(t, err)
	assert.True(t, isLocked)

	err = lock.Acquire(ctx)
	assert.Equal(t, ErrLockAcquired, err)

	err = lock.Release(ctx)
	assert.NoError(t, err)

	isLocked, err = lock.IsLocked(ctx)
	assert.NoError(t, err)
	assert.False(t, isLocked)
}
