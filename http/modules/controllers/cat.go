package controllers

import (
	"net/http"
	"sandbox/http/modules/globals"
)

func Cat(w http.ResponseWriter, r *http.Request) {
	globals.Tpl.ExecuteTemplate(w, "cat.gohtml", nil)
}
