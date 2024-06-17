package main

import (
	"log"
	"net/http"
	"os"
	graphql "procurement-app-backend/graphql"
	"procurement-app-backend/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	mux.Handle("/graphql", middleware.AuthMiddleware(srv))

	// Create a custom HTTP handler to set the Content-Type header to application/json
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Serve HTTP requests using the mux router
		mux.ServeHTTP(w, r)
	}

	// Wrap the custom HTTP handler with CORS middleware
	corsHandler := c.Handler(http.HandlerFunc(handler))

	log.Printf("Server is running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
