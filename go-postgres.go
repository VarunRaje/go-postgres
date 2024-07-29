package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
	_ "github.com/lib/pq"
)

var id int
var amount float32
var description string

func main() {
	var choice int
	db := connectPostgresDB()
	for {
		fmt.Println("Choose\n1.Insert data\n2.Read data\n3.Update data\n4.Delete data\n5.Exit")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			Insert(db)
		case 2:
			Read(db)
		case 3:
			Update(db)
		case 4:
			Delete(db)
		case 5:
			db.Close()
			os.Exit(0)
		}
	}
}

// CONNECT DB

//before connecting you have to create a database and a table in psql shell (just a base code improve these code as well as you need)

func connectPostgresDB() *sql.DB {
	var user string = ReadFromVault("appliation/go", "postgres.user")
	var host string = ReadFromVault("appliation/go", "postgres.host")
	var password string = ReadFromVault("appliation/go", "postgres.password")
	var port string = ReadFromVault("appliation/go", "postgres.port")

	connstring := "user=" + user + " dbname=bank password='" + password + "' host=" + host + " port=" + port + " sslmode=verify-full"

	db, err := sql.Open("postgres", connstring)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

// INSERT

func Insert(db *sql.DB) {
	amount = 100.00
	description = "Movie Tickets"
	insertIntoPostgres(db, amount, description)

}

func insertIntoPostgres(db *sql.DB, amount float32, description string) {
	_, err := db.Exec("INSERT INTO  cash_card(amount,description) VALUES($1,$2)", amount, description)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("value inserted")
	}
}

// READ

func Read(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM cash_card")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("id  amount    description")
		for rows.Next() {
			rows.Scan(&id, &amount, &description)
			fmt.Printf("%d - %f - %s \n", id, amount, description)
		}

	}
}

// UPDATE

func Update(db *sql.DB) {
	id = 1
	description = "Movie Tickets Updated"
	_, err := db.Exec("UPDATE cash_card SET description=$1 WHERE id=$2", description, id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data updated")
	}
}

// DELETE

func Delete(db *sql.DB) {
	id = 1
	_, err := db.Exec("DELETE FROM cash_card WHERE id=$1", id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data deleted")
	}
}

func ReadFromVault(engine string, key string) string {
	config := vault.DefaultConfig()

	config.Address = "http://127.0.0.1:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken("myroot")

	secret, err := client.KVv2("secret").Get(context.Background(), engine)
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	value, _ := secret.Data[key].(string)

	return value

}
