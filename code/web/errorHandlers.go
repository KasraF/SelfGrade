package web

import (
	"net/http"
	"io/ioutil"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	errorPage, err := ioutil.ReadFile("../templates/errors/404.html");
	
	if err != nil {
		w.WriteHeader(404)
		w.Write(errorPage)
	} else {
		logger.Error("Could not find error page at \"../templates/errors/404.html\".", err)
		ServerError(w, r)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	
	errorPage, err := ioutil.ReadFile("../templates/errors/400.html");
	
	if err != nil {
		w.Write(errorPage)
	} else {
		logger.Error("Could not find error page at \"../templates/errors/404.html\".", err)
	}
}
