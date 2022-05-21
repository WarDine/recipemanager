package main

import (
	"fmt"
	api "recipemanager/gateway/api"
	repo "recipemanager/gateway/repositories"
)

func main() {

	fmt.Println("RecipeManager-service started!")
	repo.PostgresRepo = repo.NewPostgresManager()

	api.StartServer()
	
}