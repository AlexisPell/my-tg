package db

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitScylla() {
	cluster := gocql.NewCluster("127.0.0.1") // Подставь IP своего ScyllaDB
	cluster.Keyspace = "chat_app"
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(">>> Error connecting to ScyllaDB:", err)
	}
	log.Println(">>> Successfully connected to scylla db")
}

func CloseScylla() {
	Session.Close()
}
