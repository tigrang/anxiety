package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/tigrang/anxiety"
)

//go:embed template.html
var templateString string

type templateData struct {
	Number int
}

func main() {
	// No beta blockers = panic
	anxiety.BetaBlockers = os.Getenv("ANXIETY_PANIC") != "1"

	tmpl := template.Must(template.New("test").Parse(templateString))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		// wrap errors in anxiety.Panic -- panics when BetaBlockers=false to render stack trace in browser using middleware
		err := anxiety.Panic(tmpl.Execute(&buf, templateData{Number: 10}))
		if err != nil {
			log.Printf("error rendering: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	})

	// use therapy middleware to treat the panic
	wrappedMux := anxiety.Therapy(mux)

	http.ListenAndServe(":8080", wrappedMux)
}
