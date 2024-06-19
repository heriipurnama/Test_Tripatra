package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
)

func NewGraphQLServer() *handler.Server {
	return handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: &Resolver{},
	}))
}
