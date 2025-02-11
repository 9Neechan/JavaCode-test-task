package main

import (
	"database/sql"
	"log"

	"github.com/9Neechan/JavaCode-test-task/api"
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/util"
	_ "github.com/lib/pq"
)

func main() {
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }

    conn, err := sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    store := db.NewStore(conn)
    server, err := api.NewServer(store)
    if err != nil {
        log.Fatal("cannot create server:", err)
    }

    err = server.Start(config.ServerAddress)
    if err != nil {
        log.Fatal("cannot start server:", err)
    }
}