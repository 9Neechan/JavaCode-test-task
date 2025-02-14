package main

import (
	"database/sql"
	"log"

	//api_nobalance "github.com/9Neechan/JavaCode-test-task/api_nobalance"
	api "github.com/9Neechan/JavaCode-test-task/api"
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	rabbitmq "github.com/9Neechan/JavaCode-test-task/rabbitmq"
	"github.com/9Neechan/JavaCode-test-task/redis_cache"
	"github.com/9Neechan/JavaCode-test-task/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	redis_cache.InitRedis(config.RedisAddr, config.RedisPassword, config.RedisDB)

	rabbitClient, err := rabbitmq.NewRabbitMQ(config.AmpqURL)
	if err != nil {
		log.Fatal("не удалось подключиться к RabbitMQ:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	//server, err := api_nobalance.NewServer(store)
	server, err := api.NewServer(store, rabbitClient)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
