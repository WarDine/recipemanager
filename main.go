package main

import (
	"fmt"
	api "recipemanager/gateway/api"
	repo "recipemanager/gateway/repositories"
)

func main() {

	fmt.Println("Hello motherfucker World!")
	repo.PostgresRepo = repo.NewPostgresManager()

	api.StartServer()
	
}