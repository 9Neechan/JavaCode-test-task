package api

import (
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer создает новый экземпляр сервера с заданным хранилищем и настраивает маршрутизатор.
func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store:    store,
	}

	server.setupRouter()
	return server, nil
}

// setupRouter настраивает маршрутизатор для обработки HTTP-запросов.
func (server *Server) setupRouter() {
	router := gin.Default()

	// без балансировки нагрузки
	router.POST("api/v1/wallet", server.updateWalletBalance) // http://localhost:8080/api/v1/wallet
	router.GET("api/v1/wallets/:id", server.getWallet)       // http://localhost:8080/api/v1/wallets/:id

	server.router = router
}

// Start запускает HTTP-сервер на указанном адресе.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse формирует ответ с сообщением об ошибке.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
