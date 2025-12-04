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

func queryBirdArrayByName(query string) ([]Bird, error) {
	trimmedString := strings.Trim(query, " \n\r")
	clearedString := clearString(trimmedString)
	searchString := "%" + clearedString + "%"
	log.Println("searchstring: " + searchString)
	if len(clearedString) == 0 {
		return nil, errors.New("Query string passed was empty or invalid")
	}

	rows, err := postgresConnection.Query(context.Background(), "SELECT file_path FROM bird_files WHERE name LIKE $1;", searchString)
	if err != nil {
		log.Printf("Error querying db: %v", err)
		return nil, err
	}
	defer rows.Close()

	var birdArray []Bird
	for rows.Next() {
		var filePath string
		err = rows.Scan(&filePath)
		if err != nil {
			return nil, err
		}
		var bird, err = getBird(filePath)
		if err != nil {
			return nil, err
		}
		birdArray = append(birdArray, bird)
	}

	// return getBird(birdFile)
	return birdArray, nil
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
