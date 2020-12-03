package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sandbox/http/modules/globals"
)

func Upload(w http.ResponseWriter, r *http.Request) {

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

	globals.Tpl.ExecuteTemplate(w, "upload.gohtml", data)
}
