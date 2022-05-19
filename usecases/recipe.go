package usecases

// import domain "recipemanager/domain"


type Recipe struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Calories int `json:"calories"`
	CookingTime int `json:"cookingTime"`
	Instructions string `json:"instructions,omitempty"`
	Portions int `json:"portions,omitempty"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

type Ingredient struct {
	Name string `json:"name"`
	RecipeUID int `json:"recipeUID"`
	Amount int `json:"amount"`
}

// type Ingredients struct {

// }

type Menu struct {

}

// Enforce interface
// var _ domain.Recipe = (*Recipe)(nil)


func NewRecipe(name, description, instructions string, calories, cookingTime, portions int, ingredients []Ingredient) *Recipe {
	return &Recipe {
		Name: name,
		Description: description,
		Calories: calories,
		CookingTime: cookingTime,
		Instructions: instructions,
		Portions: portions,
		Ingredients: ingredients,
	}
}

func (r *Recipe) GetRecipe() *Recipe{
	return r;
}
