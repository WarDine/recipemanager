package repositories

import (
	"recipemanager/domain"
    "fmt"
    "log"
	"os"
	"time"
	"strconv"
	"strings"
	"reflect"
	"recipemanager/usecases"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

const (
	selectMessHallInfoByIDQuery = "SELECT * FROM messhall WHERE messhalls_uid = '%s';"
	selectMessHallMenuInfoQuery = "SELECT * FROM menu WHERE menu_uid = '%s';"
)

type PostgresManager struct {
	conn *sqlx.DB
}

var PostgresRepo *PostgresManager;

// Enforce interface
var _ domain.PostgressManagerInterface = (*PostgresManager)(nil)

// return postgres objects which contains connection to database
func NewPostgresManager() *PostgresManager {

	port, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PGPORT is: ", port)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"), port, os.Getenv("PGUSER") , os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"))

	log.Println(".env: ", psqlInfo)

	// wait for database to accept connections
	time.Sleep(20 * time.Second)

	conn, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		time.Sleep(20 * time.Second)
		conn, err = sqlx.Connect("postgres", psqlInfo);
		if err != nil {
			panic(err)
		}
	}

	return &PostgresManager {
		conn: conn,
	}
}

func (pg *PostgresManager) PingConnection() {

	err := pg.conn.Ping()
	if err != nil {
		panic(err)
	} else {
		log.Println("Connection works as expected")
	}
}

func (pg *PostgresManager) DeleteConnection() {
	defer pg.conn.Close()
}

// get db tag name of a field from a generic struct
func GetDBTagName(genericStruct interface{}, structField string) string {

	tagName := "db"
	field, ok := reflect.TypeOf(genericStruct).Elem().FieldByName(structField)
	if !ok {
		log.Println("Field not found")
	}

	return string(field.Tag.Get(tagName))
}

/**
 * this generates a string like:
 * "ingredient_name, recipe_uid, amount"
 */
func createQueryFields(fields []string ) string {
	
	var queryFields strings.Builder

	for i, field := range fields {
		if i == len(fields) - 1 {
			queryFields.WriteString(field)
			break
		}

		queryFields.WriteString(field)
		queryFields.WriteString(", ")
	}

	return queryFields.String()
}

/**
 * this generates a string like:
 * ":ingredient_name, :recipe_uid, :amount"
 */
func createQueryValues(fields []string ) string {
	
	var queryFields strings.Builder

	for i, field := range fields {
		queryFields.WriteString(":")
		if i == len(fields) - 1 {
			queryFields.WriteString(field)
			break
		}

		queryFields.WriteString(field)
		queryFields.WriteString(", ")
	}

	return queryFields.String()
}


func convertfilterToString(m map[string]interface{}) {

}

/**
 * filter is a map[string]interface{}
 * it will be converted into a string an added in query
 */
 // WIP
func (pg *PostgresManager) GetFilteredIngredients(tableName string, filter string) []usecases.Ingredient {

	db := pg.conn

	if filter == "" {
		log.Print("Filter is empty; Getting all values:")
		return pg.GetAllIngredients(tableName)
	}

	ingredients := []usecases.Ingredient{}
	err := db.Select(&ingredients, "SELECT * FROM ingredient;")
	if err != nil {
		log.Println(err)
	}

	if len(ingredients) == 0 {
		log.Println("Ingredients table is empty")
	}

	return  ingredients
}

// GetAllMessHallsInfoByID
func (pg *PostgresManager) GetMessHallsInfoByID(id string) ([]usecases.MessHall, error) {

	db := pg.conn

	messHalls := []usecases.MessHall{}
	query := fmt.Sprintf(selectMessHallInfoByIDQuery, id)
	err := db.Select(&messHalls, query)

	if err != nil {
		return nil, err
	}

	return messHalls, nil
}



// GetMessHallMenuInfo for a mess hall
func (pg *PostgresManager) GetMessHallMenuInfo(messHallID string) ([]usecases.Menu, error) {

	db := pg.conn

	/* Get the messhall with this ID*/
	messHall, err := pg.GetMessHallsInfoByID(messHallID)
	if err != nil {
		return nil, err
	}

	log.Println("Messhall ", messHall)

	if len(messHall) == 0 {
		return nil, nil
	}
	/* Get the menu for this mess hall */
	messHallMenu := []usecases.Menu{}
	query := fmt.Sprintf(selectMessHallMenuInfoQuery, messHall[0].MenuUID)
	err = db.Select(&messHallMenu, query)
	if err != nil {
		return nil, err
	}

	return messHallMenu, nil
}

// getListDBMEssHallTags
func getListDBMEssHallTags(messHall *usecases.MessHall) []string {

	t := reflect.TypeOf(*messHall)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(messHall, t.Field(i).Name)
	}

	return tagFields
}

// getListDBMEssHallTags
func getListDBTagsMenu(menu *usecases.Menu) []string {

	t := reflect.TypeOf(*menu)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(menu, t.Field(i).Name)
	}

	return tagFields
}

// getListDBMEssHallAdminTags
func getListDBMEssHallAdminTags(messHall *usecases.MessHallAdmin) []string {

	t := reflect.TypeOf(*messHall)

	tagFields := make([]string, t.NumField())
	for i := range tagFields {
		tagFields[i] = GetDBTagName(messHall, t.Field(i).Name)
	}

	return tagFields
}

// AddMessHall creates a new mess hall entry
func (pg *PostgresManager) AddMessHall(messHall *usecases.MessHall, messHallAdmin *usecases.MessHallAdmin) error {

	db := pg.conn
	masshassTX := db.MustBegin()
	masshassadminsTX := db.MustBegin()

	//insert mess Hall info into table
	structTags := getListDBMEssHallTags(messHall)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", "messhall", createQueryFields(structTags), createQueryValues(structTags))
	masshassTX.NamedExec(query, &messHall)
	err := masshassTX.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	//insert mess hall admin info into table
	structTags = getListDBMEssHallAdminTags(messHallAdmin)
	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", "messhalls_admins", createQueryFields(structTags), createQueryValues(structTags))
	masshassadminsTX.NamedExec(query, &messHallAdmin)
	err = masshassadminsTX.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (pg *PostgresManager) InsertMenu(menu *usecases.Menu, messhallUID string) {

	structTags := getListDBTagsMenu(menu)
	// INSERT INTO ingredient (ingredient_uid, ingredient_name, calories) VALUES (:ingredient_uid, :ingredient_name, :calories)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", "menu", createQueryFields(structTags), createQueryValues(structTags))

	db := pg.conn

	tx := db.MustBegin()

	tx.NamedExec(query, &menu)
	err := tx.Commit()
	if err != nil {
		log.Println(err)
	}

	query = fmt.Sprintf("UPDATE messhall SET menu_uid = '%s' WHERE messhalls_uid = '%s';", menu.MenuUID, messhallUID);

	txMesshall := db.MustBegin()

	txMesshall.NamedExec(query, &menu)
	err = txMesshall.Commit()
	if err != nil {
		log.Println(err)
	}
}
