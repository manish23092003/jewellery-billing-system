package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Connect to DB directly
	dsn := "postgres://postgres:postgres@localhost:5432/jewellery_billing?sslmode=disable"
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	ctx := context.Background()
	
	// Query all unique customers from bills that don't exist in customers table
	// We'll just fetch all bills and try to insert them using the repo
	
	query := `
		SELECT organization_id, customer_name, customer_phone, grand_total 
		FROM bills 
		WHERE customer_name != ''
	`
	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()



	count := 0
	for rows.Next() {
		var orgID string
		var name, phone string
		var totalPurchases float64
		if err := rows.Scan(&orgID, &name, &phone, &totalPurchases); err != nil {
			continue
		}

		// Try to find if customer exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM customers WHERE organization_id = $1 AND name = $2 AND phone = $3)`
		err := dbPool.QueryRow(ctx, checkQuery, orgID, name, phone).Scan(&exists)
		if err == nil && !exists {
			// Insert missing customer
			insertQuery := `
				INSERT INTO customers (organization_id, name, phone, email, address, total_purchases)
				VALUES ($1, $2, $3, '', '', $4)
			`
			_, errInsert := dbPool.Exec(ctx, insertQuery, orgID, name, phone, totalPurchases)
			if errInsert == nil {
				fmt.Printf("Backfilled customer: %s (Phone: %s)\n", name, phone)
				count++
			} else {
				fmt.Printf("Failed to backfill %s: %v\n", name, errInsert)
			}
		}
	}

	fmt.Printf("Backfill complete! Added %d missing customers.\n", count)
}
