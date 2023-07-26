package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp/v3"

	_ "github.com/lib/pq"
)

type Car struct {
	ID      uint       `db:"id"`
	Details CarDetails `db:"details"`
}

type CarDetails struct {
	Make  string `json:"make"`
	Model string `json:"model"`
}

func (d *CarDetails) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, d)
	case string:
		return json.Unmarshal([]byte(src), d)
	default:
		return fmt.Errorf("unsupported type for CarDetails: %T", src)
	}
}

func (d CarDetails) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func main() {
	db, err := sqlx.Open("postgres", "user=postgres password=1234 dbname=leasing sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Create a new car
	car := Car{
		Details: CarDetails{
			Make:  "Honda",
			Model: "Civic",
		},
	}

	// Insert the car and retrieve the ID of the new car
	var id uint
	if err := db.QueryRowx(
		`INSERT INTO cars (details) VALUES ($1) RETURNING id`,
		car.Details,
	).Scan(&id); err != nil {
		panic(err)
	}
	fmt.Printf("Inserted a new car with ID: %d\n", id)

	// Retrieve the car from the database
	var carFromDB Car
	if err := db.Get(
		&carFromDB,
		`SELECT id, details FROM cars WHERE id = $1`,
		id,
	); err != nil {
		panic(err)
	}

	pp.Println(carFromDB)
}
