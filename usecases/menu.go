package usecases

type Menu struct {
	MenuUID   string `db:"menu_uid" json:"menuUID"`
	RecipeUID string `db:"recipe_uid" json:"recipeUID"`
	TimeStamp string `db:"time_stamp" json:"timeStamp"`
}