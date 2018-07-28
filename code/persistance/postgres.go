package persistance

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"GoLog"
)

var database *sql.DB = nil
var logger = GoLog.GetLogger()

// The User type queries
const (
	ex_user_createTable = "CREATE TABLE IF NOT EXISTS users (" +
		"id serial NOT NULL," +
		"name VARCHAR(255) NOT NULL," +
		"email varchar(255) NOT NULL," +
		"password bytea NOT NULL," +
		"role varchar(255) NOT NULL," +
		"PRIMARY KEY (id)" +
		");"
	qu_user_findByID    = "SELECT * FROM users WHERE users.id = $1"
	qu_user_findByEmail = "SELECT * FROM users WHERE users.email = $1"
	qu_user_save        = "INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4);"
)

// The App type queries
const (
	ex_app_createTable = "CREATE TABLE IF NOT EXISTS apps (" +
		"id serial NOT NULL," +
		"name VARCHAR(255) NOT NULL," +
		"description varchar(1024) NOT NULL," +
		"url varchar(1024) NOT NULL," +
		"iconUrl varchar(1024) NOT NULL," +
		"PRIMARY KEY (id)" +
		");"
	qu_app_findByID   = "SELECT * FROM apps WHERE apps.name = $1;"
	qu_app_findByName = "SELECT * FROM apps WHERE apps.id = $1;"
	qu_app_findAll    = "SELECT * FROM apps;"
	qu_app_save       = "INSERT INTO apps (name, description, url, iconUrl) VALUES ($1, $2, $3, $4);"
)

func InitPostgreSQL() {
	var err error
	database, err = sql.Open("postgres", "user=selfgrade password=password dbname=selfgrade sslmode=disable")

	if err != nil {
		logger.Error("Connection to PostgreSQL database failed. Exiting.", err)
		os.Exit(1)
	}

	// Create the tables
	err = createTable(ex_user_createTable)

	if err != nil {
		logger.Error("Creating database Users table failed. Exiting.", err)
		os.Exit(1)
	}

	err = createTable(ex_app_createTable)

	if err != nil {
		logger.Error("Creating database Apps table failed. Exiting.", err)
		os.Exit(1)
	}
}

func GetDatabase() *sql.DB {
	if database == nil {
		logger.Error("You need to Init()ialize the database before calling GetDatabase(). Returning nil.", nil)
	}

	return database
}

func createTable(query string) error {
	tx, err := database.Begin()

	if err != nil {
		logger.Error("Could not begin transaction to create table.", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query)

	if err != nil {
		logger.Error("Failed to create table with query \"%s\"", err, query)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func hasId(id int) bool {
	return id > 0
}

// ---------------------------------------------------------------------------------------------------------------------
// Implementations in PostgreSQL for data types.
// ---------------------------------------------------------------------------------------------------------------------

/*********
 * User
 *
 * @Todo Remove the need to pass the database to user function calls.
 ********/
func (user *User) Exists(db *sql.DB) (bool, error) {
	var rows *sql.Rows
	var err error

	if hasId(user.Id) {
		rows, err = db.Query(qu_user_findByID, user.Id)
	} else if user.Email != "" {
		rows, err = db.Query(qu_user_findByEmail, user.Email)
	}

	return rows.Next(), err
}

func (user *User) Find(db *sql.DB) (bool, error) {
	var rows *sql.Rows
	var err error
	found := false

	if hasId(user.Id) {
		rows, err = GetDatabase().Query(qu_user_findByID, user.Id)
	} else if user.Email != "" {
		rows, err = GetDatabase().Query(qu_user_findByEmail, user.Email)
	}

	if err != nil { return false, err }

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role)
		found = true
	}
	
	if err != nil { return false, err }

	if rows.Next() {
		logger.Warn("FindUser() returned multiple results for user %s", nil, user)
	}

	return found, nil
}

func (user *User) Save(db *sql.DB) error {
	_, err := GetDatabase().Exec(qu_user_save, user.Name, user.Email, user.Password, user.Role)

	if err != nil {
		logger.Error("", err)
	}
	
	return err
}

/*********
 * Apps
 ********/
func (app *App) Exists() (bool, error) {
	var rows *sql.Rows
	var err error

	if hasId(app.Id) {
		rows, err = GetDatabase().Query(qu_app_findByID, app.Id)
	} else if app.Name != "" {
		rows, err = GetDatabase().Query(qu_app_findByName, app.Name)
	}

	return rows.Next(), err
}

func (app *App) Find() (bool, error) {
	var rows *sql.Rows
	var err error
	found := false

	if hasId(app.Id) {
		rows, err = GetDatabase().Query(qu_app_findByID, app.Id)
	} else if app.Name != "" {
		rows, err = GetDatabase().Query(qu_app_findByName, app.Name)
	}

	if err != nil { return false, err }

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&app.Id, &app.Name, &app.Description, &app.Url, &app.IconUrl)
		found = true
	}
	
	if err != nil { return false, err }

	if rows.Next() {
		logger.Warn("App.Find() returned multiple results for app %s", nil, app)
	}

	return found, nil
}

func (app *App) Save() error {
	_, err := GetDatabase().Exec(qu_app_save, app.Name, app.Description, app.Url, app.IconUrl)
	
	if err != nil {
		logger.Error("", err)
	}
	
	return err
}

func FindAllApps() ([]App, error) {
	var apps []App

	rows, err := GetDatabase().Query(qu_app_findAll)

	if err != nil { return nil, err}

	defer rows.Close()

	for rows.Next() {
		var app App
		err = rows.Scan(&app.Id, &app.Name, &app.Description, &app.Url, &app.IconUrl)

		if err != nil {
			logger.Error("Failed to scan app in FindAllApp(). Skipping this app.", err)
		} else {
			apps = append(apps, app)
		}
	}

	logger.Debug("Found %u apps in call to FindAllApps().", len(apps))

	return apps, nil
}
