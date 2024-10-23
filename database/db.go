package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDataBase() {
	err := godotenv.Load("user.env")
	if err != nil {
		log.Println("Error occurs on .env file")
		fmt.Println(err)
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	database := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s database=%s password=%s sslmode=disable", host, port, user, database, pass)
	Db, err = sql.Open("postgres", psqlSetup)
	if err != nil {
		log.Printf("Error connecting the database")
		fmt.Println(err)
	} else {
		fmt.Println("Database connected successfully")
	}
}
