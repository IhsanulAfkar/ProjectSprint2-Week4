package main

import (
	"Week4/db"
	"Week4/routes"
	"context"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	db.Init(ctx)
	r :=routes.Init()
	r.Run(":8080")
}