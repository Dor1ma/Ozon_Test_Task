package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Dor1ma/Ozon_Test_Task/internal/database/storage"
	"github.com/Dor1ma/Ozon_Test_Task/pkg/graphql"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	storageType := os.Getenv("STORAGE_TYPE")

	var repo database.Repository

	if storageType == "" {
		log.Fatalf("Error: STORAGE_TYPE environment variable not set")
	}

	if storageType == "postgres" {
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			dbUser, dbPassword, dbName, dbHost, dbPort)

		db, err := pgx.Connect(context.Background(), connStr)
		if err != nil {
			log.Fatalf("failed to connect to the database: %v", err)
		}
		defer db.Close(context.Background())

		repo = storage.NewPostgreSQLRepository(db)

	} else if storageType == "in_memory" {
		repo = storage.NewInMemoryRepository()
	}

	resolver := graphql.NewResolver(repo)

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.MultipartForm{})

	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
