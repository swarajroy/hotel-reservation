package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

type TestMongoClient struct {
	Client    *mongo.Client
	Container testcontainers.Container
	Uri       string
	DbName    string
}

func NewTestMongoClient(dbName string) (*TestMongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	container, client, uri, err := createMongoContainer(ctx, dbName)
	if err != nil {
		log.Fatal("failed to setup test", err)
		return nil, err
	}

	return &TestMongoClient{
		Client:    client,
		Container: container,
		Uri:       uri,
		DbName:    dbName,
	}, nil
}

func createMongoContainer(ctx context.Context, dbName string) (testcontainers.Container, *mongo.Client, string, error) {
	var env = map[string]string{
		//"MONGO_INITDB_ROOT_USERNAME": "root",
		//"MONGO_INITDB_ROOT_PASSWORD": "pass",
		"MONGO_INITDB_DATABASE": dbName,
	}
	var port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("mongo container ready and running at port: ", p.Port())

	uri := fmt.Sprintf("mongodb://localhost:%s", p.Port())
	client, err := newMongoClient(uri)
	if err != nil {
		return container, client, "", fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, client, uri, nil
}
