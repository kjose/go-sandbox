package controllers

import (
	"custom2"
	"net/http"
	"net/url"
	"sandbox/http/modules/custom"
	"sandbox/http/modules/db"
	"sandbox/http/modules/globals"
	"sandbox/http/modules/session"
	"time"
)

func Infos(w http.ResponseWriter, r *http.Request) {
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
		User          db.User
	}{
		r.Method,
		r.URL,
		r.Host,
		r.ContentLength,
		r.Form,
		r.Header,
		cookieValue,
		session.GetSession(),
		db.DbUsers[session.Get("email")],
	}

	// Set cookie for a future visit
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: time.Now().String(),
	})

	globals.Tpl.ExecuteTemplate(w, "infos.gohtml", data)
}
