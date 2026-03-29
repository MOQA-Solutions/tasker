package main 

import (
    "database/sql"
    "log"
    "fmt"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/MOQA-Solutions/tasker/server"
)

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("CREATE DATABASE tasks")
    if err != nil {
        log.Printf("Error creating the database: %v", err)
    }
    db.Close()

	db, err = sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/tasks?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`
        CREATE TABLE api_keys (
		key TEXT PRIMARY KEY
	);
	`)

	if err != nil {
        log.Fatal(err)
    }

	_, err = db.Exec(`
        CREATE TABLE tasks (
        id    SERIAL PRIMARY KEY,
		key   TEXT NOT NULL,
        task  TEXT NOT NULL,
        state   TEXT NOT NULL
	);
    `)

	newKey, err := server.GenerateAPIKey()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("new key: %s", newKey) 

    hash := server.HashAPIKey(newKey)
    fmt.Printf("new hash: %s", hash) 

    _, err = db.Exec("INSERT INTO api_keys (key) VALUES ($1)", hash)
    if err != nil {
        log.Fatal(err)
    }

    newKey, err = server.GenerateAPIKey()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("new key: %s", newKey) 

    hash = server.HashAPIKey(newKey)
    fmt.Printf("new hash: %s", hash) 

    _, err = db.Exec("INSERT INTO api_keys (key) VALUES ($1)", hash)
    if err != nil {
        log.Fatal(err)
    }
    
	db.Close()

}

