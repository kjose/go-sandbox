package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"

	"sandbox/http/modules/session"

	"custom2"                     // vendor (marche si GOPATH configuré, se trouve dans ./vendor/custom2)
	"sandbox/http/modules/custom" // depuis le gopath
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  []byte
}

var dbUsers = make(map[string]User)

func main() {
	urls := map[string]func(w http.ResponseWriter, r *http.Request){
		"/":         infos,
		"/cat":      cat,
		"/me":       me,
		"/infos":    infos,
		"/upload":   upload,
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

var tpl *template.Template

func init() {
	// tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	tpl = template.New("")
	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gohtml") {
			tpl.ParseFiles(path)
		}

		return nil
	})

	// init user database
	bs, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.MinCost)
	dbUsers["test@abtasty.com"] = User{"Kévin", "José", "kevin.jose@abtasty.com", bs}
}

func cat(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "cat.gohtml", nil)
}

func me(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Me, me and me only")
}

func infos(w http.ResponseWriter, r *http.Request) {
	custom.Hello()
	custom2.Hello2()
	r.ParseForm()

	cookieName := "last-visit"
	var c *http.Cookie
	var cookieValue string
	c, _ = r.Cookie(cookieName)
	if c != nil {
		cookieValue = c.Value
	}

	data := struct {
		Method        string
		Url           *url.URL
		Host          string
		ContentLength int64
		Form          map[string][]string
		Header        map[string][]string
		LastVisit     string
		Session       session.Session
		User          User
	}{
		r.Method,
		r.URL,
		r.Host,
		r.ContentLength,
		r.Form,
		r.Header,
		cookieValue,
		session.GetSession(),
		dbUsers[session.Get("email")],
	}

	// Set cookie for a future visit
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: time.Now().String(),
	})

	tpl.ExecuteTemplate(w, "infos.gohtml", data)
}

func upload(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Link string
		Body string
	}

	// display uploaded image
	if r.Method == http.MethodPost {

		// body
		bs := make([]byte, r.ContentLength)
		r.Body.Read(bs)
		data.Body = string(bs)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bs))

		// open
		f, h, err := r.FormFile("q")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// fmt.Println("\nfile:", f, "\nheader:", h, "\nerr:", err)

		// read
		bs, err = ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// store on the server
		dst, err := os.Create(filepath.Join("./assets", h.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = dst.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Link = "/resources/" + h.Filename
	}

	tpl.ExecuteTemplate(w, "upload.gohtml", data)
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
		dbUsers[r.FormValue("email")] = User{
			r.FormValue("firstname"),
			r.FormValue("lastname"),
			r.FormValue("email"),
			bs,
		}

		session.Set("email", r.FormValue("email"))

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signin.gohtml", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	var message string
	if r.Method == http.MethodPost {
		if _, ok := dbUsers[r.FormValue("email")]; ok {
			session.Set("email", r.FormValue("email"))
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		message = "User not exists"
	}

	tpl.ExecuteTemplate(w, "login.gohtml", message)
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
