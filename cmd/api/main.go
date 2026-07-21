package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env Not found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Missing DATABASE_URL")
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("DB configuration error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	var databaseName string
	err = db.QueryRow(
		context.Background(),
		"SELECT current_database()",
	).Scan(&databaseName)
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}

	log.Printf("Connected to db: %s", databaseName)

	var productRowCounter int

	err = db.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM products",
	).Scan(&productRowCounter)

	if err != nil {
		log.Fatalf("Product counter error: %v", err)
	}

	log.Printf("Number of rows in products: %d", productRowCounter)
}
