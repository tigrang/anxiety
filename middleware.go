package anxiety

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"slices"
	"strconv"

	"github.com/DataDog/gostackparse"
)

func Therapy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !BetaBlockers {
			defer panicHandler(w)
		}

		next.ServeHTTP(w, r)
	})
}

func panicHandler(w http.ResponseWriter) {
	var rvr any

	if rvr = recover(); rvr == nil {
		return
	}

	stack := debug.Stack()
	goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))
	frames := goroutines[0].Stack

	// find the first panic going up the stack
	errorRefs := make([]ErrorRef, 0)
	for i := len(frames) - 1; i > 0; i-- {
		if frames[i].Func == "panic" || frames[i].Func == "github.com/tigrang/anxiety.Panic" {
			break
		}

		errorRef, err := NewErrorRef(frames[i].File, strconv.Itoa(frames[i].Line), "", frames[i].Func, "")
		if err != nil {
			io.WriteString(w, fmt.Sprintln("panic handler error: %w", err))
			return
		}

		errorRefs = append(errorRefs, errorRef)
	}

	slices.Reverse(errorRefs)

	TemplateData{
		Type:    "Runtime",
		Message: fmt.Sprintf("%v", rvr),
		Stack:   errorRefs,
	}.RenderError(w)
}
