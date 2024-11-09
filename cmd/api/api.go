package api

import (
	"github.com/LoginX/SprayDash/internal/controller"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	listenAddr string
	dbClient   *mongo.Client
}

func NewAPIServer(dbClient *mongo.Client, listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		dbClient:   dbClient,
	}
}

func (a *APIServer) Start() error {
	router := gin.Default()
	rootController := controller.NewRootController()
	rootController.RegisterRoutes(router)

	return router.Run(a.listenAddr)
}
