package web

import (
	"net/http"
	"io/ioutil"
	"strings"
)

import (
	"SelfGrade/code/security"
	"GoLog"
)

const SessionTimeout int = 30 * 60

var logger = GoLog.GetLogger()

func Init() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/", securityFilter(homeHandler))
	http.HandleFunc("/resources/", resourceHandler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Return the Login page
		loginPage, err := ioutil.ReadFile("../templates/login.html");
		
		if err != nil {
			logger.Error("Could not find login page at \"../templates/login.html\".", err)
			NotFound(w, r)
		} else {
			w.WriteHeader(200)
			w.Write(loginPage)
		}
		
	case "POST":
		err := r.ParseForm()

		if err != nil {
			// TODO Handle
			logger.Error("Parsing login form failed.", err)
			w.Header()["Location"] = []string{"/login"}
			w.WriteHeader(302)
		} else {
			form := r.PostForm

			email    := form.Get("email")
			password := form.Get("password")
			
			// TODO Handle "User does not exist" cases
			user, found := security.GetAndAuthenticateUser(email, password)
			
			if found && user.Authenticated {
				session := NewSession()
				session.user = user
				SaveSession(session)

				w.Header()["Set-Cookie"] = []string{SessionCookieName + "=" + session.Id + "; Max-Age=1800"}
				w.Header()["Location"] = []string{"/hub"}
				w.WriteHeader(302)
			} else {
				w.Header()["Location"] = []string{"/login"}
				w.WriteHeader(302)
			}
		}

	default:
		logger.Debug("Intercepted %s request to login. Interesting...", r.Method);
		NotFound(w, r)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		sessionId, sessionNotFound := r.Cookie(SessionCookieName)

		if sessionNotFound != nil {
			logger.Debug("Logout request received, but request did not contain a session cookie.. Ignoring.")
			NotFound(w, r)
			return
		}
		
		found := RemoveSession(sessionId.Value)

		if !found {
			logger.Debug("Logout request received from non-existing session. Ignoring.")
			NotFound(w, r)
			return
		}
		
		w.Header()["Set-Cookie"] = []string{SessionCookieName + "=" + sessionId.Value + "; Max-Age=0"}
		w.Header()["Location"]   = []string{"/login"}
		w.WriteHeader(302)
		
	default:
		logger.Debug("Intercepted %s request to logout. Interesting...", r.Method);
		NotFound(w, r)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()

		if err != nil {
			// TODO Handle
			logger.Error("Parsing signup form failed.", err)
			return
		}
		
		form := r.PostForm

		name     := form.Get("name")
		email    := form.Get("email")
		password := form.Get("password")
		confirm  := form.Get("confirm")
		
		if _, found := security.GetUser(email); found {
			// Handle case where user exists
			logger.Debug("Sign up request for existing email: %s", email)
			w.Header()["Location"] = []string{"/login"}
			w.WriteHeader(302)
			return
		}

		if strings.Compare(password, confirm) != 0 {
			// Handle case where password and confirm don't match
			logger.Debug("Sign up password and confirm don't match: %s vs %s", password, confirm)
			w.Header()["Location"] = []string{"/login"}
			w.WriteHeader(302)
			return
		}

		// Create user
		err = security.NewUser(name, email, password, false)
		
		if err != nil {
			logger.Error("Failed to create new user %s", err, email)
		}
		
		user, found := security.GetUser(email)
		
		if !found {
			logger.Error("Created User %s but cannot Get() it from security. Redirecting to login page.", nil, email)
			w.Header()["Location"] = []string{"/login"}
			w.WriteHeader(302)
			return
		}
		
		// User is authenicated.
		user.Authenticated = true
		
		// Create new session
		session := NewSession()
		session.user = user
		SaveSession(session)
		
		// Login the user
		w.Header()["Set-Cookie"] = []string{SessionCookieName + "=" + session.Id + "; Max-Age=1800"}
		w.Header()["Location"] = []string{"/hub"}
		w.WriteHeader(302)
		
	default:
		logger.Debug("Intercepted %s request to sign up. Interesting...", r.Method);
		NotFound(w, r)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Return the Login page
	homePage, err := ioutil.ReadFile("../templates/home.html");
	
	if err != nil {
		logger.Error("Could not find home page at \"../templates/home.html\".", err)
		NotFound(w, r)
	} else {
		w.WriteHeader(200)
		w.Write(homePage)
	}
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {

	filePath := ".." + r.URL.String()
	resource, err := ioutil.ReadFile(filePath)

	if err != nil {
		logger.Warn("Resource not found: " + filePath, err)
		NotFound(w, r)
	} else {
		filetype := filePath[strings.LastIndex(filePath,".") + 1:]

		switch filetype {
		case "css":
			w.Header().Add("Content-Type", "text/css")
		case "js":
			w.Header().Add("Content-Type", "text/javascript")
		case "svg":
			w.Header().Add("Content-Type", "image/svg+xml")
		case "map":
			fallthrough
		case "json":
			w.Header().Add("Content-Type", "app/json")
		case "woff2":
			w.Header().Add("Content-Type", "font/woff2")
		default:
			logger.Warn("Resource format not recognized for request: " + filePath, nil)
			w.Header().Add("Content-Type", "text/text")
		}

		w.WriteHeader(200)
		w.Write(resource)
	}
}
