package session

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type SessionStorage struct {
	CurrentSessionId string
	Sessions         map[string]Session
}

type Session map[string]string

var sessionStorage *SessionStorage

func init() {
	sessionStorage = &SessionStorage{
		"",
		map[string]Session{},
	}
	fmt.Println("Init session storage", sessionStorage)
}

func Init(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sid")
	sessionStorage.CurrentSessionId = ""
	if err != nil {
		Create(w, r)
		return
	}

	cval := cookie.Value
	if _, ok := sessionStorage.Sessions[cval]; !ok {
		Create(w, r)
		return
	}

	sessionStorage.CurrentSessionId = cval
	Set("sid", cval)
	fmt.Println("Init existing session id", cval)
}

func Set(name string, value string) {
	sessionStorage.Sessions[sessionStorage.CurrentSessionId][name] = value
	fmt.Println("Update session storage", sessionStorage)
}

func Get(name string) string {
	return sessionStorage.Sessions[sessionStorage.CurrentSessionId][name]
}

func GetSession() Session {
	return sessionStorage.Sessions[sessionStorage.CurrentSessionId]
}

func HasSession() bool {
	return len(sessionStorage.Sessions[sessionStorage.CurrentSessionId]) > 0
}

func Create(w http.ResponseWriter, r *http.Request) {
	cuid := uuid.New().String()
	cookie := &http.Cookie{
		Name:  "sid",
		Value: cuid,
		// Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	sessionStorage.Sessions[cuid] = Session{}
	sessionStorage.CurrentSessionId = cuid
	Set("sid", cuid)
	fmt.Println("Create new session id", cuid)
}

func Close(w http.ResponseWriter) {
	fmt.Println("Close session id", sessionStorage.CurrentSessionId)
	delete(sessionStorage.Sessions, sessionStorage.CurrentSessionId)
	sessionStorage.CurrentSessionId = ""
	cookie := &http.Cookie{
		Name:   "sid",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
