# anxiety

Render errors with stack trace directly in browser for easier development while preserving idomatic Go error handling.

![image](https://github.com/user-attachments/assets/dfb3b295-351e-41bd-bd5a-7a9bf59e0e82)

## Install
```
go get github.com/tigrang/anxiety
```

### Usage

```go
func main() {
	// No beta blockers = panic
	anxiety.BetaBlockers = false

	tmpl := template.Must(template.New("test").Parse(templateString))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		// wrap errors in anxiety.Panic
		// panics when BetaBlockers=false to render stack trace in browser using middleware
		err := anxiety.Panic(tmpl.Execute(&buf, templateData{Number: 10}))
		if err != nil {
			// regular error handling
			log.Printf("error rendering: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	})

	// use therapy middleware to treat the panic
	wrappedMux := anxiety.Therapy(mux)
	// or r.Use(anxiety.Therapy) for chi router

	http.ListenAndServe(":8080", wrappedMux)
}

```

## Run Sample

### Panic rendering off (default)

```
go run _samples/main.go
```

### Panic rendering on
```
ANXIETY_PANIC=1 go run _samples/main.go
```

Open [localhost:8080](http://localhost:8080)