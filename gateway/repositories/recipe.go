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
func (pg *PostgresManager) getListDBTagsRecipe(varName *usecases.Recipe) []string {

	t := reflect.TypeOf(*varName)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(varName, t.Field(i).Name)
	}

	return tagFields
}

func (pg *PostgresManager) InsertRecipeIntoDB(tableName string, recipe *usecases.Recipe) {

	structTags := pg.getListDBTagsRecipe(recipe)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, createQueryFields(structTags), createQueryValues(structTags))
	log.Println("Query for databse is: ", query)

	db := pg.conn
	tx := db.MustBegin()

	tx.NamedExec(query, &recipe)
	err := tx.Commit()
	if err != nil {
		log.Println(err)
	}
}

func (pg *PostgresManager) GetAllRecipes() []usecases.Recipe {

	db := pg.conn

	recipes := []usecases.Recipe{}
	err := db.Select(&recipes, "SELECT * FROM recipe;")
	if err != nil {
		log.Println(err)
	}

	if len(recipes) == 0 {
		log.Println("Recipes table is empty")
	}

	return recipes
}

func (pg *PostgresManager) GetRecipesByRecipeUID(recipeUID string) []usecases.Recipe {

	db := pg.conn
	query := fmt.Sprintf("SELECT * FROM %s WHERE recipe_uid='%s';", "recipe", recipeUID)

	recipes := []usecases.Recipe{}
	err := db.Select(&recipes, query)
	
	if err != nil {
		log.Println(err)
	}

	if len(recipes) == 0 {
		log.Println("Recipes table is empty")
	}

	return recipes
}

func (pg *PostgresManager) GetRecipesByMesshallUID(messhallUID string) []usecases.Recipe {

	db := pg.conn
	query := fmt.Sprintf("SELECT * FROM %s WHERE messhall_uid='%s';", "recipe", messhallUID)

	recipes := []usecases.Recipe{}
	err := db.Select(&recipes, query)
	
	if err != nil {
		log.Println(err)
	}

	if len(recipes) == 0 {
		log.Println("Recipes table is empty")
	}

	return recipes
}


func (pg *PostgresManager) GetRecipesCaloriesAsc() []usecases.Recipe {
	db := pg.conn

	recipes := []usecases.Recipe{}
	err := db.Select(&recipes, "SELECT * FROM recipe ORDER BY calories ASC;")
	if err != nil {
		log.Println(err)
	}

	if len(recipes) == 0 {
		log.Println("Recipes table is empty")
	}

	return recipes
}

func (pg *PostgresManager) DeleteRecipe(recipeUID string) error {

	db := pg.conn
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_uid='%s';", "recipe", recipeUID)

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return err;
	}

	return nil
}

func (pg *PostgresManager) DeleteRecipesForMesshall(messhallUID string) error {

	db := pg.conn
	query := fmt.Sprintf("DELETE FROM %s WHERE messhall_uid='%s';", "recipe", messhallUID)

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return err;
	}

	return nil
}
