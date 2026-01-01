package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Parse command line flags
	migrateOnly := flag.Bool("migrate", false, "Run migrations only and exit")
	rollbackVersion := flag.String("rollback", "", "Rollback to specific migration version")
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Get bot token from environment
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Initialize database
	// First connect to default postgres db to create the barbershop db
	tempDB, err := sql.Open("postgres", "host=localhost port=5432 user=barber password=password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = tempDB.Ping()
	if err != nil {
		log.Fatal("Cannot ping postgres database")
	}

	// Check if database already exists
	var exists int
	err = tempDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = 'barbershop'").Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking database existence: %v", err)
		return
	}
	if err == sql.ErrNoRows {
		// Database does not exist, create it
		_, err = tempDB.Exec("CREATE DATABASE barbershop")
		if err != nil {
			log.Printf("Create database error: %v", err)
		} else {
			log.Println("Database barbershop created")
		}
	} else {
		log.Println("Database barbershop already exists")
	}
	tempDB.Close()

	// Now connect to the barbershop db
	db, err = sql.Open("postgres", "host=localhost port=5432 user=barber password=password dbname=barbershop sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot ping barbershop database")
	}
	defer db.Close()

	// Handle command line flags
	if *rollbackVersion != "" {
		log.Printf("Rolling back migration: %s", *rollbackVersion)
		if err := rollbackMigration(db, *rollbackVersion); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Printf("Successfully rolled back migration: %s", *rollbackVersion)
		return
	}

	// Run database migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Exit if migrate-only flag is set
	if *migrateOnly {
		log.Println("Migrations completed successfully")
		return
	}

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}
