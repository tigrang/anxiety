package anxiety

import (
	"path/filepath"
	"regexp"
)

func ParseErrorOutput(appPath string, output string) (ErrorRef, error) {
	templeRegx := regexp.MustCompile(`file=(.+?) error=(.+?): (.+?): line (\d+), col (\d+)`)
	matches := templeRegx.FindStringSubmatch(output)

	if len(matches) > 0 {
		return NewErrorRef("templ", matches[1], matches[3], matches[4], matches[2])
	}

	goRegx := regexp.MustCompile(`(.+?):(\d+):(\d+): (.+(?:\n\t.+)?)`)
	goMatches := goRegx.FindStringSubmatch(output)

	if len(goMatches) > 0 {
		file := goMatches[1]
		fullPath, err := filepath.Abs(filepath.Join(appPath, file))
		if err != nil {
			panic(err)
		}
		return NewErrorRef(fullPath, goMatches[2], goMatches[3], "", goMatches[4])
	}

	return ErrorRef{}, nil
}
