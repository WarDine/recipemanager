package main

import (
	"fmt"
	// api "recipemanager/gateway/api"
	repo "recipemanager/gateway/repositories"
	// _ "github.com/lib/pq"
	// "database/sql"
    // "log"
	// "os"
	// "strconv"
)

func main() {
	fmt.Println("Hello motherfucker World!")

	// var pg repo.PostgresManager
	pg := repo.NewPostgresManager()
	pg.TestDatabase("ingredient")

	// fmt.Println("TRy to open Connection!")
	// port, err := strconv.Atoi(os.Getenv("PGPORT"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("PGPORT is: ", port)

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	os.Getenv("PGHOST"), port, os.Getenv("PGUSER") , os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"))
	// 	// host, port, user, password, dbname)

	// fmt.Println(".env: ", psqlInfo)

	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }


	// _ = pg.OpenConnection()
	// pg.GetData()

	// api.StartServer()
	
}