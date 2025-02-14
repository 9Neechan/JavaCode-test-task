package api

import (
	"os"
	"testing"

	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store, rabbitClient *rabbitmq.RabbitMQ) *Server {
	server, err := NewServer(store, rabbitClient)
	require.NoError(t, err)

	return server
}
