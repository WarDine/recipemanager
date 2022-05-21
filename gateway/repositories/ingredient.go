package repositories

import (
    "fmt"
    "log"
	"reflect"
	"recipemanager/usecases"
	_ "github.com/lib/pq"
)


// get list of all db tags of a struct ingredient
// call with GetListDBTags(ingredient)
func (pg *PostgresManager) getListDBTagsIngredient(ingredient *usecases.Ingredient) []string {

	t := reflect.TypeOf(*ingredient)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(ingredient, t.Field(i).Name)
	}

	return tagFields
}

// example of query
// tx.NamedExec("INSERT INTO ingredient (ingredient_name, recipe_uid, amount) VALUES (:ingredient_name, :recipe_uid, :amount)", &second_ingredient)
func (pg *PostgresManager) InsertIngredientIntoDB(tableName string, ingredient *usecases.Ingredient) {

	structTags := pg.getListDBTagsIngredient(ingredient)
	// INSERT INTO ingredient (ingredient_uid, ingredient_name, calories) VALUES (:ingredient_uid, :ingredient_name, :calories)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, createQueryFields(structTags), createQueryValues(structTags))
	db := pg.conn

	tx := db.MustBegin()

	tx.NamedExec(query, &ingredient)
	err := tx.Commit()
	if err != nil {
		log.Println(err)
	}
}

func (pg *PostgresManager) GetAllIngredients(tableName string) []usecases.Ingredient {

	db := pg.conn

	ingredients := []usecases.Ingredient{}
	err := db.Select(&ingredients, "SELECT * FROM ingredient;")
	if err != nil {
		log.Println(err)
	}

	if len(ingredients) == 0 {
		log.Println("Ingredients table is empty")
	}

	return ingredients
}

// this ingredient is unique or does not exist
func (pg *PostgresManager) GetIngredientWithUID(ingredientUID string) usecases.Ingredient {

	db := pg.conn
	query := fmt.Sprintf("SELECT * from %s WHERE ingredient_uid='%s';", "ingredient", ingredientUID)

	ingredient := usecases.Ingredient{}
	err := db.Select(&ingredient, query)
	if err != nil {
		log.Println(err)
	}

	return ingredient
}
