package main

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

var postgresConnection, _ = initPostgresConnection()

const databaseUrl string = "postgresql://postgres:1234@localhost:5432/postgres"

func closePostgresConnection(connection *pgx.Conn) {
	log.Println("Closing the connection to postgresql.")
	connection.Close(context.Background())
}

func initPostgresConnection() (*pgx.Conn, error) {
	connection, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Printf("Error connecting to PSQL: %v\n", err)
		return nil, err
	}
	log.Println("Initialized the connection to postgresql")
	return connection, nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func queryBirdTableByName(query string) (Bird, error) {
	trimmedString := strings.Trim(query, " \n\r")
	if len(trimmedString) == 0 {
		return Bird{}, errors.New("Query string passed was empty or invalid")
	}
	clearedString := clearString(trimmedString)

	var birdName string
	var birdFile string
	err := postgresConnection.QueryRow(context.Background(), "SELECT name FROM bird_files WHERE name LIKE $1;", clearedString).Scan(&birdName)
	if err != nil {
		log.Printf("Error querying db: %v", err)
		return Bird{}, errors.New("Error querying the database")
	}
	err = postgresConnection.QueryRow(context.Background(), "SELECT file_path FROM bird_files WHERE name LIKE $1;", clearedString).Scan(&birdFile)
	if err != nil {
		log.Printf("Error querying db: %v", err)
		return Bird{}, errors.New("Error querying the database")
	}

	return getBird(birdFile)
}
