package usecases

type RecipeIngredient struct {
	RecipeUID int `db:"recipe_uid" json:"recipeUID"`
	Ingredient_uid int `db:"ingredient_uid" json:"ingredientUID"`
	Amount int `db:"amount" json:"amount"`
}
