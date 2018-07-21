package persistance

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"GoLog"
)

var database *sql.DB = nil
var logger = GoLog.GetLogger()

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

func InitPostgreSQL() {
	var err error
	database, err = sql.Open("postgres", "user=selfgrade password=password dbname=selfgrade sslmode=disable")

	if err != nil {
		logger.Error("Connection to PostgreSQL database failed. Exiting.", err)
		os.Exit(1)
	}

	// Create the tables
	err = createUserTable()

	if err != nil {
		logger.Error("Creating database tables failed. Exiting.", err)
		os.Exit(1)
	}
}

func GetDatabase() *sql.DB {
	if database == nil {
		logger.Error("You need to Init()ialize the database before calling GetDatabase(). Returning nil.", nil)
	}

	return database
}

// ---------------------------------------------------------------------------------------------------------------------
// Implementations in PostgreSQL for data types.
// ---------------------------------------------------------------------------------------------------------------------

/*********
 * User
 ********/
func createUserTable() error {
	tx, err := database.Begin()

	if err != nil {
		logger.Error("Could not begin transaction to create tables.", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(ex_user_createTable)

	if err != nil {
		logger.Error("Failed to create users table.", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (user *User) Exists(db *sql.DB) (bool, error) {
	var rows *sql.Rows
	var err error

	if user.Id > -1 {
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

	if user.Id > 0 {
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
