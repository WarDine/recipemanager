package repositories

import (
	"recipemanager/domain"
    "fmt"
    "log"
	"os"
	"time"
	"strconv"
	"strings"
	"reflect"
	"recipemanager/usecases"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type PostgresManager struct {
	conn *sqlx.DB
}

var PostgresRepo *PostgresManager;

// Enforce interface
var _ domain.PostgressManagerInterface = (*PostgresManager)(nil)

// return postgres objects which contains connection to database
func NewPostgresManager() *PostgresManager {

	port, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PGPORT is: ", port)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"), port, os.Getenv("PGUSER") , os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"))

	log.Println(".env: ", psqlInfo)

	// wait for database to accept connections
	time.Sleep(20 * time.Second)

	conn, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		time.Sleep(20 * time.Second)
		conn, err = sqlx.Connect("postgres", psqlInfo);
		if err != nil {
			panic(err)
		}
	}

	return &PostgresManager {
		conn: conn,
	}
}

func (pg *PostgresManager) PingConnection() {

	err := pg.conn.Ping()
	if err != nil {
		panic(err)
	} else {
		log.Println("Connection works as expected")
	}
}

func (pg *PostgresManager) DeleteConnection() {
	defer pg.conn.Close()
}

// get db tag name of a field from a generic struct
func GetDBTagName(genericStruct interface{}, structField string) string {

	tagName := "db"
	field, ok := reflect.TypeOf(genericStruct).Elem().FieldByName(structField)
	if !ok {
		log.Println("Field not found")
	}

	return string(field.Tag.Get(tagName))
}

/**
 * this generates a string like:
 * "ingredient_name, recipe_uid, amount"
 */
func createQueryFields(fields []string ) string {
	
	var queryFields strings.Builder

	for i, field := range fields {
		if i == len(fields) - 1 {
			queryFields.WriteString(field)
			break
		}

		queryFields.WriteString(field)
		queryFields.WriteString(", ")
	}

	return queryFields.String()
}

/**
 * this generates a string like:
 * ":ingredient_name, :recipe_uid, :amount"
 */
func createQueryValues(fields []string ) string {
	
	var queryFields strings.Builder

	for i, field := range fields {
		queryFields.WriteString(":")
		if i == len(fields) - 1 {
			queryFields.WriteString(field)
			break
		}

		queryFields.WriteString(field)
		queryFields.WriteString(", ")
	}

	return queryFields.String()
}


func convertfilterToString(m map[string]interface{}) {

}

/**
 * filter is a map[string]interface{}
 * it will be converted into a string an added in query
 */
 // WIP
func (pg *PostgresManager) GetFilteredIngredients(tableName string, filter string) []usecases.Ingredient {

	db := pg.conn

	if filter == "" {
		log.Print("Filter is empty; Getting all values:")
		return pg.GetAllIngredients(tableName)
	}

	ingredients := []usecases.Ingredient{}
	err := db.Select(&ingredients, "SELECT * FROM ingredient;")
	if err != nil {
		log.Println(err)
	}

	if len(ingredients) == 0 {
		log.Println("Ingredients table is empty")
	}

	return  ingredients
}

// this function is just an example of how to use sqlx lib
/*
func (pg *PostgresManager) TestDatabase(tableName string) {

	db := pg.conn

	tx := db.MustBegin()
	var secondIngredient = &usecases.Ingredient {
		Name: "Pastarnac",
		RecipeUID: 1918273,
		Amount: 2,
	}

	tx.MustExec("INSERT INTO ingredient (ingredient_name, recipe_uid, amount) VALUES ($1, $2, $3)", "Telina", 1918273, 2)
	tx.NamedExec("INSERT INTO ingredient (ingredient_name, recipe_uid, amount) VALUES (:ingredient_name, :recipe_uid, :amount)", &secondIngredient)
	tx.MustExec("INSERT INTO ingredient (ingredient_name, recipe_uid, amount) VALUES ($1, $2, $3)", "Morcov", 1918273, 4)
	tx.Commit()

	ingredients := []usecases.Ingredient{}

	err := db.Select(&ingredients, "SELECT * FROM ingredient;")
	if err != nil {
		log.Println(err)
	}

	log.Printf("Len of ingredients: %d\n", len(ingredients))
	if len(ingredients) >= 2 {
		ingredient1, ingredient2 := ingredients[0], ingredients[1]
		log.Printf("%#v\n%#v", ingredient1, ingredient2)
	}

	pastarnac := Ingredient{}
	err = db.Get(&pastarnac, "SELECT * FROM ingredient WHERE ingredient_name=$1;", "Pastarnac")
	if err != nil {
		log.Println(err)
	}

	log.Printf("Pastarnacul: %#v\n", pastarnac)

	pg.DeleteConnection()
}

*/