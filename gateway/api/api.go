package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	// "recipemanager/domain"
	"recipemanager/usecases"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/copier"
)

const (
	HttpServerPort = ":8080"
)

var recipe usecases.Recipe

func handleGetRecipe(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(recipe)
}


func handleCreateRecipe(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var newRecipe usecases.Recipe

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newRecipe)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {

		fmt.Fprintf(os.Stdout, "Recipe: %+v\n", newRecipe)
		w.WriteHeader(http.StatusOK)
	}

	copier.Copy(&recipe, &newRecipe)
	json.NewEncoder(w).Encode(newRecipe)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewRecipeAPI() *mux.Router {
	
	var router = mux.NewRouter()
	router.Use(commonMiddleware)

	recipe.Name = "Empty"

	router.HandleFunc("/api/recipe", handleGetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/create", handleCreateRecipe).Methods("POST")

	return router
}

func StartServer() {

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})
	
	router := NewRecipeAPI()

	fmt.Printf("HTTP Server is running at http://localhost%s\n", HttpServerPort)
	log.Fatal(http.ListenAndServe(HttpServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(router)))

}