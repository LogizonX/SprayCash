package api

import (
	"github.com/LoginX/SprayDash/internal/controller"
	repo "github.com/LoginX/SprayDash/internal/repository/impls"
	service "github.com/LoginX/SprayDash/internal/service/impls"
	"github.com/LoginX/SprayDash/internal/utils"
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
	v1 := router.Group("/api/v1")
	mongoDb := a.dbClient.Database("spraydash")
	// repo
	userRepo := repo.NewUserRepositoryImpl(mongoDb)
	partyRepo := repo.NewPartyRepoImpl(mongoDb)
	// service
	cacheService := utils.NewRedisCacheService()
	mailer := utils.NewMailerService()
	codeGenerator := utils.NewCodeGeneratorService()
	userService := service.NewUserServiceImpl(userRepo, cacheService, mailer, codeGenerator)
	azureService := utils.NewAzureServiceImpl()
	partyService := service.NewPartyServiceImpl(partyRepo, azureService)
	// controller
	userController := controller.NewUserController(userService)
	partyController := controller.NewPartyController(partyService, userService)
	userController.RegisterRoutes(v1)
	partyController.RegisterPartyRoutes(v1)

	return router.Run(a.listenAddr)
}
