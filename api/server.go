package api

import (
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/rabbitmq"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewServer(store db.Store, rabbitClient *rabbitmq.RabbitMQ) (*Server, error) {
	server := &Server{
		store:    store,
		rabbitMQ: rabbitClient,
	}

	// Запускаем обработчик сообщений RabbitMQ
	for i := 0; i < 10; i++ {
		go rabbitClient.ConsumeMessages("wallet_updates", server.processUpdateWallet)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// попытка балансировки с помощью Redis RabbitMQ
	router.POST("api/v1/wallet", server.updateWalletBalanceRabbitmq) // http://localhost:8080/api/v1/wallet
	router.GET("api/v1/wallets/:id", server.getWalletRedis)          // http://localhost:8080/api/v1/wallets/:id

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
