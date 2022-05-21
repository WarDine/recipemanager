package domain

import (
	"recipemanager/usecases"
)


// PostgressManagerInterface models the behaviour of the Postgress DB manager
type PostgressManagerInterface interface {
}

type PostgresManagerCreateRecipeStruct struct {
	Recipe usecases.Recipe `json:"recipe"`
	RecipeIngredients []usecases.RecipeIngredient `json:"recipeIngredients"`
}

type AddMessHallInfoQuery struct {
	MessHallAdminNickname string
	Street                string
	City                  string
	Country               string
	Status                string
	AttendanceNumber      int
}
