package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/jewellery_billing?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	rows, err := pool.Query(context.Background(), "SELECT customer_name, customer_phone FROM bills WHERE organization_id = '856f071d-5a7b-466b-a7d9-a3ad3de64ad4'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, phone string
		if err := rows.Scan(&name, &phone); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Bill Customer: %s, Phone: %s\n", name, phone)
	}
}
