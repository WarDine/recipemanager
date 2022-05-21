package usecases

type StockIngredient struct {
	MesshallUID string `db:"messhall_uid" json:"messhallUID"`
	IngredientUID string `db:"ingredient_uid" json:"ingredientUID"`
	Amount int `db:"amount" json:"amount"`
}
