package main

import (
	"fmt"
	api "recipemanager/gateway/api"
)

func main() {
	fmt.Println("Hello World!")


	api.StartServer()
	
}