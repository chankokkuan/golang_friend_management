package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

func mongoDBConnectionString(protocol string, host string, port int) string {
	var url string
	if port == 0 {
		url = fmt.Sprintf("%s://%s/?retryWrites=true&w=majority", protocol, host)
	} else {
		url = fmt.Sprintf("%s://%s:%d", protocol, host, port)
	}

	return url
}

func NewConnection(host string, port int, database string, username string, password string) (*mongo.Database, *mongo.Client, error) {
	var uri string
	if host == "localhost" {
		uri = mongoDBConnectionString("mongodb", host, port)
	} else {
		// Ports should not be included if using srv
		uri = mongoDBConnectionString("mongodb+srv", host, 0)
	}

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.Monitor = mongotrace.NewMonitor(mongotrace.WithServiceName("friend-management-db"))
	credential := options.Credential{
		// AuthMechanism: "SCRAM-SHA-256",
		Username: username,
		Password: password,
	}
	clientOptions = clientOptions.SetAuth(credential)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, nil, err
	}

	db := client.Database(database)

	return db, client, nil
}
