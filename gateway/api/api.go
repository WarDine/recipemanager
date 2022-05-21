package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"
	"log"
	repo "recipemanager/gateway/repositories"
	"recipemanager/usecases"
	"recipemanager/domain"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/copier"
	"github.com/google/uuid"
)

const (
	HttpServerPort = ":8080"
)

func handleAddRecipe(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var recipeBlob domain.PostgresManagerCreateRecipeStruct
	var newRecipe usecases.Recipe

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&recipeBlob)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {

		fmt.Fprintf(os.Stdout, "New Recipe Blob: %+v\n", recipeBlob)
		w.WriteHeader(http.StatusOK)
	}

	// add recipe to database
	copier.Copy(&newRecipe, &recipeBlob.Recipe)
	newRecipe.RecipeUID = uuid.New().String()
	repo.PostgresRepo.InsertRecipeIntoDB("recipe", &newRecipe)

	for _, recipeIngredient := range recipeBlob.RecipeIngredients {
		recipeIngredient.RecipeUID = newRecipe.RecipeUID
		repo.PostgresRepo.InsertRecipeIngredientIntoDB("recipe_ingredients", &recipeIngredient)
	}

	json.NewEncoder(w).Encode(newRecipe)
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

	newIngredient.Ingredient_uid = uuid.New().String()
	repo.PostgresRepo.InsertIngredientIntoDB("ingredient", &newIngredient)
	json.NewEncoder(w).Encode(newIngredient)
}

// ###############
// Recipes
func handleGetRecipe(w http.ResponseWriter, r *http.Request) {
	recipes := repo.PostgresRepo.GetAllRecipes()
	json.NewEncoder(w).Encode(recipes)
}

func handleGetRecipeByCaloriesAsc(w http.ResponseWriter, r *http.Request) {
	recipes := repo.PostgresRepo.GetRecipesCaloriesAsc()
	json.NewEncoder(w).Encode(recipes)
}

func handleGetRecipeByRecipeUID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	recipeUID := params["recipeUID"]

	ingredients := repo.PostgresRepo.GetRecipesByRecipeUID(recipeUID)
	json.NewEncoder(w).Encode(ingredients)
}

func handleGetRecipeByMesshallUID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messhallUID := params["messhallUID"]

	ingredients := repo.PostgresRepo.GetRecipesByMesshallUID(messhallUID)
	json.NewEncoder(w).Encode(ingredients)
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello World")
}

func handleDeleteRecipe(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	recipeUID := params["recipeUID"]
	
	err := repo.PostgresRepo.DeleteRecipe(recipeUID)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("Cannot delete recipe!")
	}
	err = repo.PostgresRepo.DeleteRecipeIngredientsForRecipe(recipeUID)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("Cannot delete recipe Ingredients")
	}
}

func handleDeleteRecipesForMesshall(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messhallUID := params["messhallUID"]

	err := repo.PostgresRepo.DeleteRecipeIngredientsForMesshall(messhallUID)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("Cannot delete recipe ingredients for this messhall!")
	}

	err = repo.PostgresRepo.DeleteRecipesForMesshall(messhallUID)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("Cannot delete recipes for this messhall!")
	}

}

// ###############
// RECIPES_INGREDIENTS
func handleGetRecipeIngredientsByUID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	recipeUID := params["recipeUID"]

	recipeIngredients := repo.PostgresRepo.GetIngredientsForRecipe(recipeUID)
	json.NewEncoder(w).Encode(recipeIngredients)
}

func handleGetRecipeIngredients(w http.ResponseWriter, r *http.Request) {
	recipeIngredients := repo.PostgresRepo.GetRecipeIngredientsTable()
	json.NewEncoder(w).Encode(recipeIngredients)
}

// ###############
// INGREDIENTS
func handleGetIngredient(w http.ResponseWriter, r *http.Request) {
	ingredients := repo.PostgresRepo.GetAllIngredients("ingredient")
	json.NewEncoder(w).Encode(ingredients)
}

func handleGetIngredientWithUID(w http.ResponseWriter, r *http.Request) {
	ingredients := repo.PostgresRepo.GetAllIngredients("ingredient")
	json.NewEncoder(w).Encode(ingredients)
}


func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handleGetStock(w http.ResponseWriter, r *http.Request) {
	stock := repo.PostgresRepo.GetStock()
	json.NewEncoder(w).Encode(stock)
}

func handleGetStockMesshall(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messhallUID := params["messhallUID"]

	stock := repo.PostgresRepo.GetStockOfMesshall(messhallUID)
	json.NewEncoder(w).Encode(stock)
}

func handleGetStockIngredient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messhallUID := params["messhallUID"]
	ingredientUID := params["ingredientUID"]

	stock := repo.PostgresRepo.GetStockOfAlimentFromMesshall(messhallUID, ingredientUID)
	json.NewEncoder(w).Encode(stock)
}

func GetMessHallMenu(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messhallUID := params["messhallUID"]

	messHallMenuEntries, err := repo.PostgresRepo.GetMessHallMenuInfo(messhallUID)
	if err != nil {
		return
	}

	log.Println("Messhall menu entries: ", messHallMenuEntries)
	menuRecipes := []usecases.Recipe{}
	for i, menuEntry := range messHallMenuEntries {

		log.Printf("Menu entry number %d : %v\n", i, messHallMenuEntries)

		recipe := repo.PostgresRepo.GetRecipesByRecipeUID(menuEntry.RecipeUID)
		log.Println("Recipe: ", recipe)

		menuRecipes = append(menuRecipes, recipe[0])
	}

	json.NewEncoder(w).Encode(menuRecipes)
}

func handleAddStock(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var newStock usecases.StockIngredient
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newStock)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprintf(os.Stdout, "Stock: %+v\n", newStock)
		w.WriteHeader(http.StatusOK)
	}

	repo.PostgresRepo.InsertStockIntoDB(&newStock)
	json.NewEncoder(w).Encode(newStock)
}

func GenerateShoppingList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messhallUID := params["messhallUID"]

	
	shoppingList := make(map[string]int)
	stockMap := make(map[string]int)
	recipeMap := make(map[string]int)

	stock := repo.PostgresRepo.GetStockOfMesshall(messhallUID)
	menu, _ := repo.PostgresRepo.GetMessHallMenuInfo(messhallUID)

	if menu == nil {
		log.Println("Menu does not exist")
		json.NewEncoder(w).Encode("Menu does not exist")
		return
	}

	// create map with stock
	for _, stockEntry := range stock {
		stockMap[stockEntry.IngredientUID] = stockEntry.Amount
	}

	log.Println("Current Stock: ", stockMap)

	recipeIngredients := []usecases.RecipeIngredient{}

	for _, menuEntry := range menu {
		recipeIngredients = append(recipeIngredients, repo.PostgresRepo.GetIngredientsForRecipe(menuEntry.RecipeUID)...)
	}
	
	// create map with all ingredients we need
	for _, ingredientsEntry := range recipeIngredients {
		recipeMap[ingredientsEntry.Ingredient_uid] += ingredientsEntry.Amount
	}

	log.Println("Ingredients for recipes: ", recipeMap)

	// create shopping list
	for ingredientUID, amount := range recipeMap {

		stockIngredientAmount, exist := stockMap[ingredientUID]
		if !exist {
			log.Println("!!")
			shoppingList[ingredientUID] = amount
			continue
		}

		if  stockIngredientAmount < amount {
			log.Println("??")
			shoppingList[ingredientUID] = amount - stockIngredientAmount
		} 
	}

	log.Println("Shopping List: ", shoppingList)

	json.NewEncoder(w).Encode(shoppingList)
}


//API
// AddMessHall add a new mess hall and its admin
// func AddMessHall(request *restful.Request, response *restful.Response) {
func AddMessHall(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var queryBody domain.AddMessHallInfoQuery
	// var recipeBlob domain.PostgresManagerCreateRecipeStruct

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&queryBody)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprintf(os.Stdout, "messhall info body: %+v\n", queryBody)
		w.WriteHeader(http.StatusOK)
	}

	/* Populate new mess hall struct */
	messHallUID := uuid.New().String()
	newMessHall := usecases.MessHall{
		MessHallUID:      messHallUID,
		Street:           queryBody.Street,
		City:             queryBody.City,
		Country:          queryBody.Country,
		Status:           queryBody.Status,
		AttendanceNumber: queryBody.AttendanceNumber,
	}

	/* Populat new mess hall admin struct */
	newMessHallAdmin := usecases.MessHallAdmin{
		MessHallAdminUID: uuid.New().String(),
		Nickname:         queryBody.MessHallAdminNickname,
		MessHallUID:      messHallUID,
	}

	/* Add mess hall info and its admin info to repository */
	err = repo.PostgresRepo.AddMessHall(&newMessHall, &newMessHallAdmin)
	if err != nil {
		json.NewEncoder(w).Encode("ERROR: failed to add new mess hall")
		if err != nil {
			return
		}
		return
	}

	err = json.NewEncoder(w).Encode(newMessHall)
	if err != nil {
		return
	}
}


func AddMenu(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	params := mux.Vars(r)
	messhallUID := params["messhallUID"]

	var newMenu usecases.Menu
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newMenu)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprintf(os.Stdout, "Menu: %+v\n", newMenu)
		w.WriteHeader(http.StatusOK)
	}

	newMenu.MenuUID = uuid.New().String()
	newMenu.TimeStamp = time.Now().Format("01-02-2006")
	repo.PostgresRepo.InsertMenu(&newMenu, messhallUID)
	json.NewEncoder(w).Encode(newMenu)
}

func NewRecipeAPI() *mux.Router {
	
	var router = mux.NewRouter()
	router.Use(commonMiddleware)

	router.HandleFunc("/", handleHelloWorld).Methods("GET")

	router.HandleFunc("/api/recipe", handleGetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/{recipeUID}", handleGetRecipeByRecipeUID).Methods("GET")
	router.HandleFunc("/api/recipe/get-by-messhall/{messhallUID}", handleGetRecipeByMesshallUID).Methods("GET")
	router.HandleFunc("/api/recipe/get-asc-by-calories", handleGetRecipeByCaloriesAsc).Methods("GET")
	router.HandleFunc("/api/recipe/add", handleAddRecipe).Methods("POST")
	router.HandleFunc("/api/recipe/delete/{recipeUID}", handleDeleteRecipe).Methods("DELETE")
	router.HandleFunc("/api/recipe/delete-by-messhall/{messhallUID}", handleDeleteRecipesForMesshall).Methods("DELETE")

	router.HandleFunc("/api/ingredient", handleGetIngredient).Methods("GET")
	router.HandleFunc("/api/ingredient/{recipeUID}", handleGetIngredientWithUID).Methods("GET")
	router.HandleFunc("/api/ingredient/add", handleAddIngredient).Methods("POST")
	router.HandleFunc("/api/ingredient/generate/{messhallUID}", GenerateShoppingList).Methods("GET")

	router.HandleFunc("/api/recipe-ingredients", handleGetRecipeIngredients).Methods("GET")
	router.HandleFunc("/api/recipe-ingredients/{recipeUID}", handleGetRecipeIngredientsByUID).Methods("GET")

	router.HandleFunc("/api/stock", handleGetStock).Methods("GET")
	router.HandleFunc("/api/stock-of-messhall/{messhallUID}", handleGetStockMesshall).Methods("GET")
	router.HandleFunc("/api/stock-of-ingredient/{messhallUID}&{ingredientUID}", handleGetStockIngredient).Methods("GET")
	router.HandleFunc("/api/stock/add", handleAddStock).Methods("POST")

	router.HandleFunc("/api/messhall/add", AddMessHall).Methods("POST")

	router.HandleFunc("/api/menu/get-meshall-menu/{messhallUID}", GetMessHallMenu).Methods("GET")
	router.HandleFunc("/api/menu/add/{messhallUID}", AddMenu).Methods("POST")


	// router.HandleFunc("/api/recipe-ingredients/add", handleAddIngredient).Methods("POST")
	// router.HandleFunc("/api/recipe-ingredients/add", handleAddIngredient).Methods("POST")

	return router
}

func StartServer() {

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})
	
	router := NewRecipeAPI()

	fmt.Printf("HTTP Server is running at http://localhost%s\n", HttpServerPort)
	log.Fatal(http.ListenAndServe(HttpServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(router)))

}

// INSERT INTO menu(menu_uid, recipe_uid, time_stamp) VALUES ('1', '2260f2b4-db58-4ef2-93e9-08c265d01f3f', '2016-06-22 19:10:25-07');
// INSERT INTO menu(menu_uid, recipe_uid, time_stamp) VALUES ('1', 'f30642fb-08a8-4671-b6df-d48603d4a06a', '2016-06-22 19:10:25-07');

// UPDATE messhall SET menu_uid = '1'  WHERE messhalls_uid = 'ba7c94ce-537a-4c64-a961-6928ae0ea252';
