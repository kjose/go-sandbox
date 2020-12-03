package controllers

import (
	"fmt"
	"net/http"
)

func Me(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Me, me and me only")
}
