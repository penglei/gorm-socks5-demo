package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/net/proxy"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(dialFunc pgconn.DialFunc, dsn string) (*gorm.DB, error) {
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	connConfig.DialFunc = dialFunc

	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	pgGormDialector := postgres.New(
		postgres.Config{
			Conn: db,
		},
	)
	return gorm.Open(pgGormDialector)

}

func NewSocks5PgxDialer(addr string) (pgconn.DialFunc, error) {
	s5proxy, err := proxy.SOCKS5("tcp", addr, nil, nil)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return s5proxy.Dial(network, addr)
	}, nil
}

func main() {

	pgxDialer, err := NewSocks5PgxDialer("127.0.0.1:18080")
	if err != nil {
		panic(err)
	}

	dsn := "postgres://postgres@localhost:5432/postgres?sslmode=disable"

	type Book struct {
		Id     int    `json:"id" gorm:"primaryKey"`
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	book := Book{
		Title:  time.Now().String(),
		Author: "fixpoint",
	}

	DB, err := NewGormDB(pgxDialer, dsn)
	if result := DB.Create(&book); result.Error != nil {
		log.Fatalln(result.Error)
	}
}
