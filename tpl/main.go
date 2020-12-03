package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// Magic Note

type MagicNote struct {
	Note int
}

func (n MagicNote) DoMagic() int {
	return 20
}

// Main

var tpl *template.Template

var fm = template.FuncMap{
	"uc":      strings.ToUpper,
	"ft":      firstThree,
	"noteFmt": noteFmt,
}

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*"))
}

func firstThree(s string) string {
	s = strings.TrimSpace(s)
	s = s[:3]
	return s
}

func noteFmt(note int) string {
	return fmt.Sprint(note) + "/20"
}

func main() {
	nf, err := os.Create("index.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer nf.Close()

	err = tpl.ExecuteTemplate(nf, "tpl.gohtml", nil)

	err = tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", struct {
		FirstName string
		LastName  string
		Notes     []MagicNote
	}{
		"Kévin",
		"José",
		[]MagicNote{MagicNote{10}, MagicNote{11}, MagicNote{13}},
	})
	if err != nil {
		log.Fatalln(err)
	}
}
