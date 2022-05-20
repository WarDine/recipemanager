package usecases

import (
	// "reflect"
)

type Ingredient struct {
	Ingredient_uid int	`db:"ingredient_uid" json:"ingredientUID,omitempty"`
	Name string `db:"ingredient_name" json:"name"`
	Calories int `db:"calories" json:"calories"`
}


// // get list of all db tags of a struct ingredient
// // call with GetListDBTags(ingredient)
// func (ing *Ingredient) GetListDBTags() []string {

// 	t := reflect.TypeOf(*ing)

// 	tagFields := make([]string, t.NumField())
// 	for i := range tagFields {
// 		tagFields[i] = GetDBTagName(ing, t.Field(i).Name)
// 	}

// 	return tagFields
// }

// struct used in postgres
	// type Ingredient struct {
	// 	Ingredient_uid int `db:"ingredient_uid" json:"ingredientUID"`
	// 	Name string `db:"ingredient_name" json:"name"`
	// 	RecipeUID int `db:"recipe_uid" json:"recipeUID"`
	// 	Amount int `db:"amount" json:"amount"`
	// }