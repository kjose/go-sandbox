package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	"sandbox/http/modules/controllers"
	"sandbox/http/modules/db"
	"sandbox/http/modules/globals"
	"sandbox/http/modules/session"
	// vendor (marche si GOPATH configuré, se trouve dans ./vendor/custom2)
	// depuis le gopath
)

func main() {
	urls := map[string]func(w http.ResponseWriter, r *http.Request){
		"/":         controllers.Infos,
		"/cat":      controllers.Cat,
		"/me":       controllers.Me,
		"/infos":    controllers.Infos,
		"/upload":   controllers.Upload,
		"/redirect": redirect,
		"/expire":   expire,
		"/login":    login,
		"/signin":   signin,
		"/logout":   logout,
	}

	for url, f := range urls {
		http.Handle(url, wrapper(http.HandlerFunc(f)))
	}

	// http.HandleFunc("/cat.jpg", catImg)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("assets"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	fmt.Println("Server started ...")

	http.ListenAndServe(":8080", nil)
}

func wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Init session
		session.Init(w, r)
		next.ServeHTTP(w, r)
	})
}

func init() {
	// tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	globals.Tpl = template.New("")
	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gohtml") {
			globals.Tpl.ParseFiles(path)
		}

		return nil
	})

	// init user database
	bs, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.MinCost)
	db.DbUsers["test@abtasty.com"] = db.User{"Kévin", "José", "kevin.jose@abtasty.com", bs}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func expire(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "last-visit",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func signin(w http.ResponseWriter, r *http.Request) {
	if isConnected() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		bs, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.DbUsers[r.FormValue("email")] = db.User{
			r.FormValue("firstname"),
			r.FormValue("lastname"),
			r.FormValue("email"),
			bs,
		}

		session.Set("email", r.FormValue("email"))

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	globals.Tpl.ExecuteTemplate(w, "signin.gohtml", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	if isConnected() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var message string
	if r.Method == http.MethodPost {
		// test user
		if u, ok := db.DbUsers[r.FormValue("email")]; ok {
			// test password
			err := bcrypt.CompareHashAndPassword(u.Password, []byte(r.FormValue("password")))
			if err == nil {
				session.Set("email", r.FormValue("email"))
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			message = "Wrong password"
		} else {
			message = "User not exists"
		}

	}

	globals.Tpl.ExecuteTemplate(w, "login.gohtml", message)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session.Close(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func isConnected() bool {
	return session.Get("email") != ""
}

// func catImg(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "resources/cat.jpg")
// }
