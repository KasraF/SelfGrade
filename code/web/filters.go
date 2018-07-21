package web

import (
	"net/http"
	"strings"
)

func securityFilter(h http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Security pages
		if strings.Compare(r.URL.String(), "/login.html") == 0 {
			h(w, r) // call original
			return
		}

		// Redirects to "/login" if new/anonymous session
		sessionId, sessionNotFound := r.Cookie(SessionCookieName)

		if sessionNotFound == nil {

			// Session cookie exists. See if we know about it
			session, found := FindSession(sessionId.Value)

			if found {
				// We know about the session. See if authenticated

				if session.Authenticated() {
					// Hurray! serve the page!
					h(w, r)
				} else {
					// Anonymous session. Redirect to login.
					w.Header()["Location"] = []string{"/login"}
					w.WriteHeader(302)
				}
				
			} else {
				// We don't know about the session create a new one and ask user to log in.
				session = NewSession()
				w.Header()["Set-Cookie"] = []string{SessionCookieName + "=" + session.Id + "; Max-Age=1800"}
				w.Header()["Location"] = []string{"/login"}
				w.WriteHeader(302)
			}
		} else {
			// No session cookie. A new user... Exciting!
			session := NewSession()
			w.Header()["Set-Cookie"] = []string{SessionCookieName + "=" + session.Id + "; Max-Age=1800"}
			w.Header()["Location"] = []string{"/login"}
			w.WriteHeader(302)
		}
	})
}

