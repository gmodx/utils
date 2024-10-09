package appfeatures

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAppFeatures_Init(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"32777:27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	host, err := mongoC.Host(ctx)
	assert.NoError(t, err)
	port, err := mongoC.MappedPort(ctx, "27017")
	assert.NoError(t, err)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port.Port()))
	client, err := mongo.Connect(ctx, clientOptions)
	assert.NoError(t, err)

	appName := "testApp"
	af := New(client, "testDB", "testCollection", appName)

	err = af.Init(ctx)
	assert.NoError(t, err)

	var entity appFeaturesEntity
	filter := bson.D{{Key: "app", Value: appName}}
	collection := client.Database("testDB").Collection("testCollection")
	err = collection.FindOne(ctx, filter).Decode(&entity)
	assert.NoError(t, err)
	assert.Equal(t, appName, entity.App)
}

func TestAppFeatures_FeaturesUpgrade(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	host, err := mongoC.Host(ctx)
	assert.NoError(t, err)
	port, err := mongoC.MappedPort(ctx, "27017")
	assert.NoError(t, err)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port.Port()))
	client, err := mongo.Connect(ctx, clientOptions)
	assert.NoError(t, err)

	appName := "testApp"
	af := New(client, "testDB", "testCollection", appName)

	err = af.Init(ctx)
	assert.NoError(t, err)

	mockFeature := MockFeature{}
	err = af.EnableFeaturesAndUpgrade(ctx, map[string]AppFeature{
		"feature1": &mockFeature,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, mockFeature.Sum)
	assert.Equal(t, &mockFeature, af.Features["feature1"])
}

type MockFeature struct {
	Sum int
}

func (m *MockFeature) Upgrade() error {
	m.Sum++
	return nil
}
