package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	repo "recipemanager/gateway/repositories"
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

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello World")
}

func handleGetIngredient(w http.ResponseWriter, r *http.Request) {
	ingredients := repo.PostgresRepo.GetAllIngredients("ingredient")
	json.NewEncoder(w).Encode(ingredients)
}

func handleAddIngredient(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var newIngredient usecases.Ingredient

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newIngredient)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprintf(os.Stdout, "Ingredient: %+v\n", newIngredient)
		w.WriteHeader(http.StatusOK)
	}

	repo.PostgresRepo.InsertIngredientIntoDB("ingredient", &newIngredient)
	json.NewEncoder(w).Encode(newIngredient)
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

	router.HandleFunc("/", handleHelloWorld).Methods("GET")
	router.HandleFunc("/api/recipe", handleGetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/create", handleCreateRecipe).Methods("POST")

	router.HandleFunc("/api/ingredient", handleGetIngredient).Methods("GET")
	router.HandleFunc("/api/ingredient/add", handleAddIngredient).Methods("POST")

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