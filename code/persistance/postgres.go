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

// The Module type queries
const (
	ex_module_createTable = "CREATE TABLE IF NOT EXISTS modules (" +
		"id serial NOT NULL," +
		"name VARCHAR(255) NOT NULL," +
		"description varchar(1024) NOT NULL," +
		"url varchar(1024) NOT NULL," +
		"iconUrl varchar(1024) NOT NULL," +
		"PRIMARY KEY (id)" +
		");"
	qu_module_findByID   = "SELECT * FROM modules WHERE modules.name = $1;"
	qu_module_findByName = "SELECT * FROM modules WHERE modules.id = $1;"
	qu_module_findAll    = "SELECT * FROM modules;"
	qu_module_save       = "INSERT INTO modules (name, description, url, iconUrl) VALUES ($1, $2, $3, $4);"
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

	err = createTable(ex_module_createTable)

	if err != nil {
		logger.Error("Creating database Modules table failed. Exiting.", err)
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
 * Modules
 ********/
func (module *Module) Exists() (bool, error) {
	var rows *sql.Rows
	var err error

	if hasId(module.Id) {
		rows, err = GetDatabase().Query(qu_module_findByID, module.Id)
	} else if module.Name != "" {
		rows, err = GetDatabase().Query(qu_module_findByName, module.Name)
	}

	return rows.Next(), err
}

func (module *Module) Find() (bool, error) {
	var rows *sql.Rows
	var err error
	found := false

	if hasId(module.Id) {
		rows, err = GetDatabase().Query(qu_module_findByID, module.Id)
	} else if module.Name != "" {
		rows, err = GetDatabase().Query(qu_module_findByName, module.Name)
	}

	if err != nil { return false, err }

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&module.Id, &module.Name, &module.Description, &module.Url, &module.IconUrl)
		found = true
	}
	
	if err != nil { return false, err }

	if rows.Next() {
		logger.Warn("Module.Find() returned multiple results for module %s", nil, module)
	}

	return found, nil
}

func (module *Module) Save() error {
	_, err := GetDatabase().Exec(qu_module_save, module.Name, module.Description, module.Url, module.IconUrl)
	
	if err != nil {
		logger.Error("", err)
	}
	
	return err
}

func FindAllModules() ([]Module, error) {
	var modules []Module

	rows, err := GetDatabase().Query(qu_module_findAll)

	if err != nil { return nil, err}

	defer rows.Close()

	for rows.Next() {
		var module Module
		err = rows.Scan(&module.Id, &module.Name, &module.Description, &module.Url, &module.IconUrl)

		if err != nil {
			logger.Error("Failed to scan module in FindAllModule(). Skipping this module.", err)
		} else {
			modules = append(modules, module)
		}
	}

	logger.Debug("Found %u modules in call to FindAllModules().", len(modules))

	return modules, nil
}
