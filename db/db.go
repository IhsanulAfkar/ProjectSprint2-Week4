package db

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)
var db *sqlx.DB
var err error
func Init(ctx context.Context) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
	)
	fmt.Println(connStr)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	db, err = sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		panic(err.Error())
	}
	db.DB.SetMaxOpenConns(100)
	db.DB.SetConnMaxIdleTime(time.Second  * 5)
	db.DB.SetConnMaxLifetime(time.Hour)
	// db.DB.SetMaxIdleConns(10)
	// Migrations
	// migrate -database "postgresql://root:root@localhost:5432/belimang?sslmode=disable" -path db/migrations up
}
func CreateConn() *sqlx.DB{
	return db
}