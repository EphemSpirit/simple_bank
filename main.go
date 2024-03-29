package main

import (
	"database/sql"
	"log"

	"github.com/EphemSpirit/simple_bank/api"
	db "github.com/EphemSpirit/simple_bank/db/sqlc"
	"github.com/EphemSpirit/simple_bank/util"
	_ "github.com/lib/pq"
)

func main() {	
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't start server")
	}
}