package security

import (
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	"fmt"
	"GoLog"
)

import (
	"SelfGrade/code/utils"
	"SelfGrade/code/persistance"
)

/** 
 * This struct only stores the security details, without the account info.
 */
type User struct {
	Email string
	Role string
	Authenticated bool
}

const ADMIN_ROLE = "ADMIN"
const USER_ROLE  = "USER"

var logger *GoLog.Logger
var db     *sql.DB

func Init() {
	logger = GoLog.GetLogger()
	db     = persistance.GetDatabase()
}

func NewUser(name string, email string, password string, admin bool) error {
	user := persistance.User{Email: email}
	exists, err := user.Find(db)

	if err == nil {
		if !exists {
			passhash, err := hashPassword(password)

			if err == nil {
				user.Name     = name
				user.Password = passhash

				if admin {
					user.Role = ADMIN_ROLE
				} else {
					user.Role = USER_ROLE
				}
					
				user.Save(db)
			}
		} else {
			err = utils.Error{Message: fmt.Sprintf("An account with email \"%s\" already exists.", user.Email)}
		}
	}

	return err
}

func GetUser(email string) (User, bool) {
	var user User
	database_user := persistance.User{Email: email}
	found, err := database_user.Find(db)
	
	if err != nil {
		logger.Error("Trying to Find() user with email \"%s\" threw an error.", err, email)
		found = false
	} else if found {
		user.Email = database_user.Email
		user.Authenticated = false
		user.Role = database_user.Role
	}
	
	return user, found
}

func GetAndAuthenticateUser(email string, password string) (User, bool) {
	user, found := GetUser(email)

	if found {
		user.Authenticate(password)
	}

	return user, found
}

func hashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	if err != nil {
		// TODO Why would this ever fail?
		logger.Error("Failed to hash the given password: \"%s\"", err, password)
		return nil, err
	}

	return hash, nil
}

func (user *User) Authenticate(password string) bool {
	var authenticated bool
	database_user := persistance.User{Email: user.Email}
	found, err := database_user.Find(db)

	if err != nil {
		logger.Error("Trying to Find() user with email \"%s\" threw an error.", err, user.Email)
		authenticated = false
	} else if !found {
		logger.Debug("Called Authenticate() for non-existing user %s. Ignoring.", user.Email)
		authenticated = false
	} else {
		user.Authenticated = bcrypt.CompareHashAndPassword(database_user.Password, []byte(password)) == nil
		authenticated = user.Authenticated
	}

	return authenticated
}
