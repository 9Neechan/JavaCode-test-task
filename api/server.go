package api

import (
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  db.Store // *
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("api/v1/wallet", server.update)     // http://localhost:8080/users
	router.GET("api/v1/wallets/:id", server.get) // http://localhost:8080/accounts/214 api/v1/wallets/{WALLET_UUID}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
