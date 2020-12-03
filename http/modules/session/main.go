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
	var cuid string
	if err != nil {
		cuid = uuid.New().String()
		cookie = &http.Cookie{
			Name:  "sid",
			Value: cuid,
			// Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		fmt.Println("Init current session id", cuid)
	}

	sessionStorage.CurrentSessionId = cuid
}

func Set(name string, value string) {
	if _, ok := sessionStorage.Sessions[sessionStorage.CurrentSessionId]; !ok {
		sessionStorage.Sessions[sessionStorage.CurrentSessionId] = Session{}
	}

	sessionStorage.Sessions[sessionStorage.CurrentSessionId][name] = value
	fmt.Println("Update session storage", sessionStorage)
}

func Get(name string) string {
	return sessionStorage.Sessions[sessionStorage.CurrentSessionId][name]
}

func GetSession() Session {
	return sessionStorage.Sessions[sessionStorage.CurrentSessionId]
}
