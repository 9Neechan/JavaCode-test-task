package api

import (
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store // *
	router *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	//router.POST("api/v1/wallet", server.updateWalletBalance) // http://localhost:8080/api/v1/wallet
	//router.GET("api/v1/wallets/:id", server.getWallet)       // http://localhost:8080/api/v1/wallets/:id

	router.POST("api/v1/wallet", server.updateWalletBalanceRedis) // http://localhost:8080/api/v1/wallet
	router.GET("api/v1/wallets/:id", server.getWalletRedis)       // http://localhost:8080/api/v1/wallets/:id

	server.router = router
}

/*
POST
{
  "amount": 1,
  "wallet_uuid": "7ff05ab9-80d5-40d0-8037-7133da806e49",
  "operation_type": "DEPOSIT"
}
*/

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
