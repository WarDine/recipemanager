package usecases

type Ingredient struct {
	Ingredient_uid string `db:"ingredient_uid" json:"ingredientUID"`
	Name string `db:"ingredient_name" json:"name"`
	Calories int `db:"calories" json:"calories"`
}
