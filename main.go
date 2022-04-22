package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

func InsertUser(u *User, conn *pgx.Conn) {
	// Executing SQL query for insertion
	if _, err := conn.Exec(context.Background(), "INSERT INTO USERS(USERNAME) VALUES($1)", u.UserName); err != nil {
		// Handling error, if occur
		fmt.Println("Unable to insert due to: ", err)
		return
	}
	log.Println("Insertion Succesfull")
}

func GetAllUsers(conn *pgx.Conn) {
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
	// variable to store username
	var username string

	// Executing query for single row
	if err := conn.QueryRow(context.Background(), "SELECT username from users WHERE ID=$1", id).Scan(&username); err != nil {
		fmt.Println("Error occur while finding user: ", err)
		return
	}
	log.Printf("User with id=%v is %v\n", id, username)
}

func createTable(conn *pgx.Conn) error {
	query := `
			CREATE TABLE IF NOT EXISTS USERS(
		        ID          SERIAL   PRIMARY KEY,
		        USERNAME    VARCHAR(20) NOT NULL UNIQUE
		    );
			`
	if _, err := conn.Query(context.Background(), query); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func main() {
	// Open up our database connection.
	time_init := time.Now()

	conn, err := pgx.Connect(context.Background(), "postgres://admin:admin@localhost:5490/test")

	if err != nil {
		fmt.Println(err)
	}

	// defer the close till after the main function has finished
	// executing
	defer conn.Close(context.Background())
	// time execution
	defer log.Println("Execution time: ", time.Since(time_init))

	//Creating temporary user object.
	tmpUser := User{
		UserName: "Ricardo Arturo"}

	InsertUser(&tmpUser, conn)
	GetAnUser(7, conn)
	GetAllUsers(conn)

}
