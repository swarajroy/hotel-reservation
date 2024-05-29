package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/armory-io/go-commons/awaitility"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	var port nat.Port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port.Port()},
			WaitingFor:   wait.ForListeningPort(port),
			Env:          env,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}
	//defer container.Terminate(ctx) // not required as we are terminating the conatiner in the lifecycle method of testify when the container
	//gets used in integration tests

	var p nat.Port
	err = awaitility.AwaitDefault(func() bool {
		p, err = container.MappedPort(ctx, port)
		return err == nil
	})

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
