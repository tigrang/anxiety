package anxiety

import (
	"bufio"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed error.html
var errorHtml string

var errorTemplate *template.Template = template.Must(template.New("error").Parse(errorHtml))

var BetaBlockers bool = true

type ErrorRef struct {
	File        string
	Line        string
	Col         string
	Description string
	Func        string
	Code        string
	LineStart   int
}

type TemplateData struct {
	Type     string
	Message  string
	Stack    []ErrorRef
	ProxyUrl string
}

func (d TemplateData) RenderError(w http.ResponseWriter) {
	if err := errorTemplate.Execute(w, d); err != nil {
		io.WriteString(w, fmt.Sprintf("render error: %s", err))
	}
}

func Panic(err error) error {
	if err != nil && !BetaBlockers {
		panic(err)
	}
	return err
}

func NewErrorRef(file string, line string, col string, function string, desc string) (ErrorRef, error) {
	funcParts := strings.Split(function, ".")

	ref := ErrorRef{
		Description: desc,
		File:        file,
		Func:        funcParts[len(funcParts)-1],
		Line:        line,
		Col:         col,
	}

	lineNum, err := strconv.Atoi(line)
	if err != nil {
		panic(err)
	}

	ref.LineStart = lineNum - 10
	if ref.LineStart < 1 {
		ref.LineStart = 1
	}

	code, err := readLinesInRange(ref.File, ref.LineStart, lineNum+10)
	if err != nil {
		return ref, err
	}

	ref.Code = strings.Join(code, "\n")

	return ref, nil
}

func readLinesInRange(filePath string, startLine, endLine int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		if lineNumber >= startLine && lineNumber <= endLine {
			lines = append(lines, scanner.Text())
		}
		if lineNumber > endLine {
			break
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
