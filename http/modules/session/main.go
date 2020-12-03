package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SessionStorage struct {
	CurrentSessionId string
	Sessions         map[string]Session
}

type Session struct {
	ExpireOn time.Time
	Data     map[string]string
}

var sessionStorage *SessionStorage

var expireAfterSeconds int = 1800

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

	// Expire
	if time.Now().After(sessionStorage.Sessions[cval].ExpireOn) {
		fmt.Println("Session expired")
		Close(w)
	}
}

func Set(name string, value string) {
	sessionStorage.Sessions[sessionStorage.CurrentSessionId].Data[name] = value
	fmt.Println("Update session storage", sessionStorage)
}

func Get(name string) string {
	return sessionStorage.Sessions[sessionStorage.CurrentSessionId].Data[name]
}

func GetSession() Session {
	return sessionStorage.Sessions[sessionStorage.CurrentSessionId]
}

func HasSession() bool {
	return len(sessionStorage.Sessions[sessionStorage.CurrentSessionId].Data) > 0
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

	sessionStorage.Sessions[cuid] = Session{
		time.Now().Add(time.Duration(expireAfterSeconds) * time.Second),
		map[string]string{},
	}
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
