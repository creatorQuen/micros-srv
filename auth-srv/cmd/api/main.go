package main

import (
	"auth-srv/repository"
	"os"
	"time"

	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models repository.Models
}

func main() {
	log.Println("Starting auth service")

	conn := connectToDataBase()
	if conn == nil {
		log.Panic("Not connect to database.")
	}

	app := Config{
		DB:     conn,
		Models: repository.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDataBase() *sql.DB {
	dns := os.Getenv("DNS")

	for {
		connection, err := openDB(dns)
		if err != nil {
			log.Println("Databe not ready ...")
			counts++
		} else {
			log.Println("Conneted to db.")
			return connection
		}

		if counts > 5 {
			log.Println(err)
			return nil
		}

		log.Println("Retry to connecting before 2 seconds ...")
		time.Sleep(2 * time.Second)

		continue
	}
}
