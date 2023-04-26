package main

import (
	"fmt"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func main() {
	// Open a connection to the PostgreSQL database
	// Connect to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Close the database connection
	if err := db.Close(); err != nil {
		panic(err)
	}
	defer db.Close()

	// Connect to the database using gorm
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.DB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "myapp_",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	defer gormDB.Close()

	// Create a new migration instance
	m := gormigrate.New(gormDB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20220101000000",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					gorm.Model
					Name  string `gorm:"size:255"`
					Email string `gorm:"uniqueIndex"`
				}
				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "20220102000000",
			Migrate: func(tx *gorm.DB) error {
				type Product struct {
					gorm.Model
					Name        string
					Description string
					Price       float64
				}
				return tx.AutoMigrate(&Product{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("products")
			},
		},
		{
			ID: "20220103000000",
			Migrate: func(tx *gorm.DB) error {
				type Order struct {
					gorm.Model
					UserID uint
					Total  float64
				}
				if err := tx.AutoMigrate(&Order{}); err != nil {
					return err
				}
				if err := tx.Model(&Order{}).AddForeignKey("user_id", "myapp_users(id)", "CASCADE", "CASCADE").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropForeignKey("orders", "user_id"); err != nil {
					return err
				}
				return tx.Migrator().DropTable("orders")
			},
		},
	})

	// Run the migrations
	err = m.Migrate()
	if err != nil {
		panic(err)
	}

	// Query for all users and print their names
	var users []User
	err = db.Select(&users, "SELECT * FROM users")
	if err != nil {
		panic(err)
	}
}