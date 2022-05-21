package repositories

import (
    "fmt"
    "log"
	"reflect"
	"recipemanager/usecases"
	_ "github.com/lib/pq"
)


func (pg *PostgresManager) getListDBTagsStockIngredient(varName *usecases.StockIngredient) []string {

	t := reflect.TypeOf(*varName)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(varName, t.Field(i).Name)
	}

	return tagFields
}

func (pg *PostgresManager) InsertStockIntoDB(stock *usecases.StockIngredient) {

	structTags := pg.getListDBTagsStockIngredient(stock)
	// INSERT INTO ingredient (ingredient_uid, ingredient_name, calories) VALUES (:ingredient_uid, :ingredient_name, :calories)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", "stock_ingredient", createQueryFields(structTags), createQueryValues(structTags))
	db := pg.conn

	tx := db.MustBegin()

	tx.NamedExec(query, &stock)
	err := tx.Commit()
	if err != nil {
		log.Println(err)
	}
}

func (pg *PostgresManager) GetStock() []usecases.StockIngredient {

	db := pg.conn

	stocks := []usecases.StockIngredient{}
	err := db.Select(&stocks, "SELECT * FROM stock_ingredient;")
	if err != nil {
		log.Println(err)
	}

	if len(stocks) == 0 {
		log.Println("We do not have any stocks")
	}

	return stocks
}

func (pg *PostgresManager) GetStockOfMesshall(messhallUID string) []usecases.StockIngredient {

	db := pg.conn
	query := fmt.Sprintf("SELECT * from %s WHERE messhall_uid='%s';", "stock_ingredient", messhallUID)

	stock := []usecases.StockIngredient{}
	err := db.Select(&stock, query)
	if err != nil {
		log.Println(err)
	}

	return stock
}

func (pg *PostgresManager) GetStockOfAlimentFromMesshall(messhallUID string, ingredientUID string) []usecases.StockIngredient {

	db := pg.conn
	query := fmt.Sprintf("SELECT * from %s WHERE messhall_uid='%s' AND ingredient_uid='%s';", "stock_ingredient", messhallUID, ingredientUID, )

	stock := []usecases.StockIngredient{}
	err := db.Select(&stock, query)
	if err != nil {
		log.Println(err)
	}

	return stock
}

