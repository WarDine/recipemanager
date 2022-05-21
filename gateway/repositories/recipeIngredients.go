package repositories

import (
    "fmt"
    "log"
	"reflect"
	"recipemanager/usecases"
	_ "github.com/lib/pq"
)


func (pg *PostgresManager) getListDBTagsRecipeIngredients(varName *usecases.RecipeIngredient) []string {

	t := reflect.TypeOf(*varName)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(varName, t.Field(i).Name)
	}

	return tagFields
}

func (pg *PostgresManager) InsertRecipeIngredientIntoDB(tableName string, recipeIngredient *usecases.RecipeIngredient) {

	structTags := pg.getListDBTagsRecipeIngredients(recipeIngredient)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, createQueryFields(structTags), createQueryValues(structTags))
	log.Println("Query for databse is: ", query)
	db := pg.conn

	tx := db.MustBegin()

	tx.NamedExec(query, &recipeIngredient)
	err := tx.Commit()
	if err != nil {
		log.Println(err)
	}
}

// returns ingredient from recipe_ingredients table
func (pg *PostgresManager) GetIngredientsForRecipe(recipeUID string) []usecases.RecipeIngredient {

	query := fmt.Sprintf("SELECT * from %s WHERE recipe_uid='%s';", "recipe_ingredients", recipeUID)
	db := pg.conn
	recipeIngredients := []usecases.RecipeIngredient{}

	err := db.Select(&recipeIngredients, query)
	if err != nil {
		log.Println(err)
	}

	if len(recipeIngredients) == 0 {
		log.Println("This Recipe does not have ingredients!")
	}

	return recipeIngredients
}

// returns ingredient details from ingredient table
func (pg *PostgresManager) GetIngredientsDetailsForRecipe(recipeUID string) []usecases.Ingredient {

	recipeIngredients := pg.GetIngredientsForRecipe(recipeUID)
	ingredients := []usecases.Ingredient{}

	for _, recipeIngredient := range recipeIngredients {
		ingredients = append(ingredients, pg.GetIngredientWithUID(recipeIngredient.Ingredient_uid))
	}

	if len(ingredients) == 0 {
		log.Println("This Recipe does not have ingredients!")
	}

	return ingredients
}

func (pg *PostgresManager) GetRecipeIngredientsTable() []usecases.RecipeIngredient {

	query := fmt.Sprintf("SELECT * from %s;", "recipe_ingredients")
	db := pg.conn
	recipeIngredients := []usecases.RecipeIngredient{}

	err := db.Select(&recipeIngredients, query)
	if err != nil {
		log.Println(err)
	}

	if len(recipeIngredients) == 0 {
		log.Println("RecipeIngredient table does not have any values!")
	}

	return recipeIngredients
}

func (pg *PostgresManager) DeleteRecipeIngredientsForRecipe(recipeUID string) error {

	db := pg.conn
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_uid='%s';", "recipe_ingredients", recipeUID)

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return err;
	}

	return nil
}


func (pg *PostgresManager) DeleteRecipeIngredientsForMesshall(messhallUID string) error {

	db := pg.conn

	recipes := pg.GetRecipesByMesshallUID(messhallUID)
	for _, recipe := range recipes {
		query := fmt.Sprintf("DELETE FROM %s WHERE recipe_uid='%s';", "recipe_ingredients", recipe.RecipeUID)
		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
			return err;
		}
	}


	return nil
}

