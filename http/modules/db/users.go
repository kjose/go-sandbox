package db

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  []byte
}

var DbUsers = make(map[string]User)
