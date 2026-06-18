package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Connect to default "postgres" db
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Create database
	_, err = conn.Exec(context.Background(), "CREATE DATABASE jewellery_billing;")
	if err != nil {
		fmt.Printf("Database might already exist or error: %v\n", err)
	} else {
		fmt.Println("Database 'jewellery_billing' created successfully!")
	}
	os.Exit(0)
}
