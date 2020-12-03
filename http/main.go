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
)

func main() {
	http.Handle("/", http.HandlerFunc(infos))
	http.Handle("/cat", http.HandlerFunc(cat))
	http.Handle("/me", http.HandlerFunc(me))
	http.HandleFunc("/infos", infos)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/expire", expire)
	// http.HandleFunc("/cat.jpg", catImg)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("assets"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	fmt.Println("Server started ...")
	http.ListenAndServe(":8080", nil)
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
}

func cat(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "cat.gohtml", nil)
}

func me(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Me, me and me only")
}

func infos(w http.ResponseWriter, r *http.Request) {
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
	}{
		r.Method,
		r.URL,
		r.Host,
		r.ContentLength,
		r.Form,
		r.Header,
		cookieValue,
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
		}
		defer f.Close()

		// fmt.Println("\nfile:", f, "\nheader:", h, "\nerr:", err)

		// read
		bs, err = ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

// func catImg(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "resources/cat.jpg")
// }
