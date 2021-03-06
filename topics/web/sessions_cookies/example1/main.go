// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to use sessions in your web app.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// sessionName contains the session key.
const sessionName = "ultimate-web-session"

// store represents the session store for the app. Don't hard code the key like this.
var store = sessions.NewCookieStore([]byte("something-very-secret"))

// htmlNoSession contains the document we will use we
// have a request has not submited state yet.
var htmlNoSession = `
<html>
    <form action="/save" method="POST">
        <label>What is your name?</label><br>
        <input type="text" name="myName" placeholder="Name goes here">
        <input type="submit" value="Submit">
    </form>
</html>`

// htmlWithSession contains the document we will use when we have a request
// that has already submited state. We're just using printf for the example
// instead of a full template.
var htmlWithSession = `
<html>
    <h1>Hello %s!</h1>
</html>`

// App loads all of our routes
func App() http.Handler {

	// Create a new mux which will process the requests.
	m := http.NewServeMux()

	// Load the two routes.
	m.HandleFunc("/", homeHandler)
	m.HandleFunc("/save", saveHandler)

	return m
}

// homeHandler provides support for the home page route.
func homeHandler(res http.ResponseWriter, req *http.Request) {

	// Look for any session related to this request.
	session, err := store.Get(req, sessionName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there is no saved state in the sesssion for this request,
	// provide the document to request it.
	name := session.Values["name"]
	if name == nil {
		fmt.Fprint(res, htmlNoSession)
		return
	}

	// There is saved state so return a document with
	// the saved state.
	fmt.Fprintf(res, htmlWithSession, name)
}

// saveHandler provides support for save route.
func saveHandler(res http.ResponseWriter, req *http.Request) {

	// Look for any session related to this request.
	session, err := store.Get(req, sessionName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the raw query from the URL and update req.Form.
	req.ParseForm()

	// Locate the myName form value.
	name := req.FormValue("myName")

	// Save this value inside the session store.
	session.Values["name"] = name

	// You must call Save before writing to the response
	// or returning from the handler.
	session.Save(req, res)

	// Print our template including the name.
	fmt.Fprintf(res, htmlWithSession, name)
}

func main() {

	// Start the http server to handle requests
	log.Fatal(http.ListenAndServe(":3000", App()))
}
