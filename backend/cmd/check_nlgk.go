package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/jewellery_billing?sslmode=disable"
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	var name, phone string
	err = dbPool.QueryRow(context.Background(), "SELECT customer_name, customer_phone FROM bills WHERE customer_name = 'nlgk' LIMIT 1").Scan(&name, &phone)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Bill found: Name=%s, Phone='%s'\n", name, phone)
		
		var exists bool
		dbPool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM customers WHERE name = 'nlgk')").Scan(&exists)
		fmt.Printf("Customer nlgk exists in customers table? %v\n", exists)
	}
}
