package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LoginX/SprayDash/cmd/api"
	"github.com/LoginX/SprayDash/config"
	"github.com/LoginX/SprayDash/db"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// get the listenAdd
	addr := config.GetEnv("LISTEN_ADDR", "localhost:2100")
	dbUrl := config.GetEnv("MONGO_URL", "mongodb://localhost:27017")
	fmt.Println(dbUrl)
	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dbClient, err := db.ConnectDb(ctx, dbUrl)
	if err != nil {
		panic(err)
	}

	// instantiate new apiserver
	apiServer := api.NewAPIServer(dbClient, addr)

	if err := apiServer.Start(); err != nil {
		panic(err)
	}

}
