package usecases

type RecipeIngredient struct {
	RecipeUID string `db:"recipe_uid" json:"recipeUID"`
	Ingredient_uid string `db:"ingredient_uid" json:"ingredientUID"`
	Amount int `db:"amount" json:"amount"`
}
