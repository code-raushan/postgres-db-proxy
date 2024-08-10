package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    proxyAddr := ":5432"
    listener, err := net.Listen("tcp", proxyAddr)
    if err != nil {
        log.Fatal("failed to create tcp server")
    }
    defer listener.Close()

    fmt.Printf("Proxy listening on %s\n", proxyAddr)

    for {
        clientConn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        go handleConnection(clientConn)
    }
}

func handleConnection(clientConn net.Conn) {
    defer clientConn.Close()

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbHost := os.Getenv("PG_DATABASE_HOST")
    dbPort := os.Getenv("PG_DATABASE_PORT")

    // Establish a new TCP connection to the actual PostgreSQL server
    dbConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", dbHost, dbPort))
    if err != nil {
        log.Printf("Failed to connect to database: %v", err)
        return
    }
    defer dbConn.Close()

    // Proxy data between client and database
    go func() {
        _, err := io.Copy(dbConn, clientConn)
        if err != nil {
            log.Printf("Error copying data from client to DB: %v", err)
        }
    }()

    _, err = io.Copy(clientConn, dbConn)
    if err != nil {
        log.Printf("Error copying data from DB to client: %v", err)
    }
}