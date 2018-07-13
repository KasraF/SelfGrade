package security

import (
	"golang.org/x/crypto/bcrypt"
	"GoLog"
)

type User struct {
	Username string
	Password []byte
	Authenticated bool
}

var users = make(map[string]User)
var logger = GoLog.GetLogger()

func NewUser(username string) {
	users[username] = User{Username: username, Authenticated: false}
}

func GetUser(username string) (User, bool) {
	user, ok := users[username]
	return user, ok
}

func UpdatePassword(username string, password string) {
	user, ok := users[username]

	if ok {
		var err error
		user.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		if err != nil {
			// TODO Why would this ever fail?
			logger.Error("Failed to hash the given password: \"%s\"", err, password)
		} else {
			users[username] = user
		}
	} else {
		logger.Debug("Called UpdatePassword() for non-existing user %s. Ignoring.", username)
	}
}

func Authenticate(username string, password string) bool {
	user, ok := users[username]

	if !ok {
		logger.Debug("Called Authenticate() for non-existing user %s. Ignoring.", username)
		return false
	} else {
		user.Authenticated = bcrypt.CompareHashAndPassword(user.Password, []byte(password)) == nil
		users[username] = user
		return user.Authenticated
	}
}
