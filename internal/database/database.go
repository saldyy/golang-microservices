package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Db struct {
	conn *pgx.Conn
}

func New(url string) *Db {

	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
    return &Db{
      conn: nil,
    }
	}

	defer conn.Close(context.Background())

  fmt.Println("inhere")
  fmt.Println(conn)
  fmt.Println(err)
	return &Db{
		conn: conn,
	}
}

func (db *Db) Check() (string, error) {
  err := db.conn.Ping(context.Background())
  if err != nil {
    return "down", err
  }

  return "up", nil
}
