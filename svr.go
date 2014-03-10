package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var port = flag.Int("p", 9500, "the port")

var store = sessions.NewCookieStore([]byte("secret santa"))
var templates *template.Template

const sessionName = "amzn"

type order struct {
	Amount string
	Desc   string
}

type user struct {
	Name  string
	Email string
}

func init() {
	// For cookie serialization
	gob.Register(&order{})
	gob.Register(&user{})

	// Init the templates
	templates = template.Must(template.New("app").ParseGlob("web/tmpl/*.html"))
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		amount := req.FormValue("amount")
		desc := req.FormValue("desc")

		// Get the session and save the values
		session, _ := getSession(req)
		session.Values["order"] = &order{Amount: amount, Desc: desc}
		session.Values["callbackURL"] = req.FormValue("callbackURL")
		session.Values["id"] = req.FormValue("id")

		// Save
		session.Save(req, w)

		// Show the login template
		templates.ExecuteTemplate(w, "login.html", req.Host)
	case "POST":
		session, _ := getSession(req)
		session.Values["user"] = &user{Name: "Francisco J Carriedo", Email: req.FormValue("email")}
		session.Save(req, w)

		// Redirect to the order details
		http.Redirect(w, req, "/orders", http.StatusFound)
	}
}

func orderHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		session, _ := getSession(req)
		params := map[string]interface{}{
			"order": session.Values["order"],
			"user":  session.Values["user"],
			"id":    session.Values["id"],
		}

		templates.ExecuteTemplate(w, "order-confirm.html", params)
	case "POST":
		// Get the ID
		id := string(req.FormValue("id"))

		// Post back to callbackURL with id
		session, _ := getSession(req)
		callbackURL := session.Values["callbackURL"]
		if res, err := http.PostForm(callbackURL.(string), url.Values{"id": []string{id}, "authToken": []string{"a4B388"}}); err != nil {
			http.Error(w, "Problem confirming your order. You will not be charged.", res.StatusCode)
			return
		}

		fmt.Fprint(w, "Order confirmed. Your items should be delivered in any moment.")
	}
}

func getSession(req *http.Request) (*sessions.Session, error) {
	session, err := store.Get(req, sessionName)
	if err != nil {
		return nil, err
	}

	session.Options = &sessions.Options{Path: "/", MaxAge: 90, HttpOnly: true}

	return session, nil
}

func main() {
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/auth", authHandler)
	router.HandleFunc("/orders", orderHandler)

	http.Handle("/", router)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Println("Amazon payment app listening on addr", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
