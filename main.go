package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

func InsertUser(u *User, conn *pgx.Conn) {
	log.Println("inserto usuarios")
	// Executing SQL query for insertion
	if _, err := conn.Exec(context.Background(), "INSERT INTO USERS(USERNAME) VALUES($1)", u.UserName); err != nil {
		// Handling error, if occur
		fmt.Println("Unable to insert due to: ", err)
		return
	}
	log.Println("Insertion Succesfull")
}

func GetAllUsers(conn *pgx.Conn) {
	log.Println("obtengo usuarios")
	if rows, err := conn.Query(context.Background(), "SELECT * FROM USERS"); err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var u User

			rows.Scan(&u.ID, &u.UserName)
			log.Printf("%+v\n", u)
		}
		if rows.Err() != nil {
			// if any error occurred while reading rows.
			fmt.Println("Error will reading user table: ", err)
			return
		}
	}
}

func GetAnUser(id int, conn *pgx.Conn) {
	log.Println("obtengo 1 usuario")
	// variable to store username
	var username string

	// Executing query for single row
	if err := conn.QueryRow(context.Background(), "SELECT username from users WHERE ID=$1", id).Scan(&username); err != nil {
		fmt.Println("Error occur while finding user: ", err)
		return
	}
	log.Printf("User with id=%v is %v\n", id, username)
}

// create table initialize
func CreateTable(conn *pgx.Conn) error {
	log.Println("creo tabla")
	query := `
			CREATE TABLE IF NOT EXISTS user_4(
		        ID          SERIAL   PRIMARY KEY,
		        USERNAME    VARCHAR(20) NOT NULL UNIQUE
		    );
			`
	if rows, err := conn.Query(context.Background(), query); err != nil {
		log.Fatal(err)
		return err
	} else {
		defer rows.Close()
	}
	return nil
}

func main() {
	// Open up our database connection.
	time_init := time.Now()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config := Config{
		ServerName: os.Getenv("ServerName"),
		User:       os.Getenv("User"),
		Password:   os.Getenv("Password"),
		DB:         os.Getenv("DB"),
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.User, config.Password, config.ServerName, config.DB)

	conn, err := pgx.Connect(context.Background(), connectionString)

	if err != nil {
		fmt.Println(err)
	}

	// defer the close till after the main function has finished
	defer conn.Close(context.Background())

	//Creating temporary user object.
	tmpUser := User{
		UserName: "UserNameTest"}

	CreateTable(conn)
	GetAllUsers(conn)
	GetAnUser(5, conn)
	InsertUser(&tmpUser, conn)

	log.Println("Execution time: ", time.Since(time_init))


}
