package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDb(ctx context.Context, dbUrl string) (*mongo.Client, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, errors.New("context has no deadline")
	}
	if time.Now().After(deadline) {
		return nil, errors.New("context deadline exceeded")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// ping database
	if err = client.Database("admin").RunCommand(ctx, map[string]interface{}{"ping": 1}).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Database connection successfully")
	return client, nil

}
