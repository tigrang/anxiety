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


## Proxy

Anxiety also comes with an `anxious` proxy that renders compile errors in the browser the same way.

### Installation

```
go install github.com/tigrang/anxiety/cmd/anxious@latest
```

### Usage

First, create a build script within your app's directory. The default build script name is `build`.

Next, start the `anxious` proxy and configure it to your app's path and address.

```
anxious --app /path/to/your/app --proxy http://localhost:3001
```

Finally, update your `air` config with the followng:

```
[proxy]
  enabled = true
  proxy_port = 3001
  app_port = 3000

[build]
  cmd = "anxious --notify"
```