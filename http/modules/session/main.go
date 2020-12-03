package session

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Session map[string]string

var sessions map[string]Session

var currentUuid string

func Init(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		uuid := uuid.New().String()
		cookie = &http.Cookie{
			Name:  "sid",
			Value: uuid,
			// Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		sessions[uuid] = Session{}

		fmt.Println("Init session", cookie)
	}

	currentUuid = cookie.Value
}

func Set(name string, value string) {
	sessions[currentUuid][name] = value
}

func Get(name string) Session {
	return sessions[currentUuid]
}
