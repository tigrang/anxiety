package anxiety

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
)

type Handler struct {
	appPath     string
	buildFail   bool
	buildOutput []byte
	buildPath   string
	cmd         string
	errorRef    ErrorRef
	proxy       *httputil.ReverseProxy
	proxyUrl    string
}

func NewProxy(appPath string, cmd string, proxyUrl string, buildPath string) *Handler {
	remote, err := url.Parse(proxyUrl)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	return &Handler{
		appPath:   appPath,
		buildPath: buildPath,
		cmd:       cmd,
		proxy:     proxy,
		proxyUrl:  proxyUrl,
	}
}

func (h *Handler) render(w http.ResponseWriter) {
	TemplateData{
		Type:     "Compile",
		Stack:    []ErrorRef{h.errorRef},
		Message:  h.errorRef.Description,
		ProxyUrl: h.proxyUrl,
	}.RenderError(w)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving request %s\n", r.URL)
	if r.URL.Path == h.buildPath {
		success, err := h.Build()

		if err != nil {
			w.Write([]byte(fmt.Sprintln("anxious build error: %w", err)))
			return
		}

		if !success {
			h.render(w)
		}

		return
	}

	if h.buildFail {
		h.render(w)
		return
	}

	h.proxy.ServeHTTP(w, r)
}

func (h *Handler) Build() (bool, error) {
	fmt.Println("Executing build step...")

	cmd := exec.Command(h.cmd)
	cmd.Dir = h.appPath

	if out, err := cmd.CombinedOutput(); err != nil {
		h.buildFail = true
		h.buildOutput = out

		if errorRef, err := ParseErrorOutput(h.appPath, string(out)); err != nil {
			return false, err
		} else {
			h.errorRef = errorRef
		}

		return false, err
	}

	h.buildFail = false
	return true, nil
}
