package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/elokanugrah/go-financing-btpns/internal/config"
	"github.com/elokanugrah/go-financing-btpns/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db := database.NewConnection(cfg)
	defer db.Close()

	if err := seedTenors(db); err != nil {
		log.Fatalf("FATAL: Failed to seed tenors: %v", err)
	}

	if err := seedUsers(db); err != nil {
		log.Fatalf("FATAL: Failed to seed users: %v", err)
	}

	log.Println("Seeding process completed successfully!")
}

// seedTenors clears the tenors table and inserts new dummy data.
func seedTenors(db *sql.DB) error {
	log.Println("Clearing tenors table...")
	_, err := db.Exec(`TRUNCATE TABLE tenors RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("error truncating tenors table: %w", err)
	}

	stmt, err := db.Prepare(`INSERT INTO tenors (tenor_value) VALUES ($1)`)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	tenors := []int{6, 12, 18, 24, 30, 36}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	// Insert 6 tenors
	for _, tenor := range tenors {
		if _, err := tx.Stmt(stmt).Exec(tenor); err != nil {
			// If any insert fails, roll back the entire transaction
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("error executing insert and rolling back transaction: %w, %w", err, rbErr)
			}
			return fmt.Errorf("error executing insert: %w", err)
		}
	}

	// Commit the transaction.
	log.Println("Committing transaction...")
	return tx.Commit()
}

// seedUsers clears the users table and inserts dummy users.
func seedUsers(db *sql.DB) error {
	log.Println("Clearing users table...")
	_, err := db.Exec(`TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("error truncating users table: %w", err)
	}

	stmt, err := db.Prepare(`INSERT INTO users (name, phone) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	users := []struct {
		Name  string
		Phone string
	}{
		{"Budi Santoso", "081234567890"},
		{"Siti Aminah", "081987654321"},
		{"Andi Wijaya", "081223344556"},
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	for _, user := range users {
		if _, err := tx.Stmt(stmt).Exec(user.Name, user.Phone); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("error executing insert and rolling back transaction: %w, %w", err, rbErr)
			}
			return fmt.Errorf("error executing insert: %w", err)
		}
	}

	log.Println("Committing transaction for users...")
	return tx.Commit()
}
