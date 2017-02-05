package pattern

import (
	"os"
	"strings"
	"bytes"
	"regexp"
)

func NormalizeDirSep(path string) (string) {
	return strings.Replace(path, "\\", "/", -1)
}

func ParsePathPattern(sourcePattern string) (string, string) {
	path := sourcePattern;
	pattern := ""
	pathExists := false
	for {
		if _, err := os.Stat(path); err == nil {
			pattern = strings.Replace(sourcePattern, path, "", 1)
			if len(pattern) > 0 {
				pattern = pattern[1:]
			}
			pathExists = true
			break
		}

		lastSlashIndex := strings.LastIndex(NormalizeDirSep(path), "/")
		if lastSlashIndex == -1 {
			break
		}
		path = path[0:lastSlashIndex]
	}

	if ! pathExists {
		path = ""
		pattern = sourcePattern
	}
	return path, pattern
}

func GlobToRegex(glob string) (string) {
	var buffer bytes.Buffer
	r := strings.NewReader(glob)

	escape := false
	braceOpen := 0
	for {
		r, _, err := r.ReadRune()
		if (err != nil) {
			break;
		}

		if (escape) {
			buffer.WriteRune(r)
			escape = false
			continue
		}

		if (r == '\\') {
			buffer.WriteRune(r)
			escape = true;
			continue
		}

		if (r == '*') {
			buffer.WriteString(".*")
			continue
		}

		if (r == '.') {
			buffer.WriteString("\\.")
			continue
		}

		if (r == '{') {
			buffer.WriteString("(")
			braceOpen++
			continue
		}

		if (r == '}') {
			buffer.WriteString(")")
			braceOpen--
			continue
		}

		if (r == ',' && braceOpen > 0) {
			buffer.WriteString("|")
			continue
		}

		buffer.WriteRune(r)
	}

	return buffer.String()
}

func BuildMatchList(sourcePattern *regexp.Regexp, subject string)([]string) {
	list := make([]string, 0)
	normalizedPath := NormalizeDirSep(subject)
	sourcePattern.ReplaceAllStringFunc(normalizedPath, func(m string) string {
		parts := sourcePattern.FindStringSubmatch(m)
		i := 1
		for range parts[1:] {
			list = append(list, parts[i])
			i++

		}
		return m
	})
	return list
}

func CompileNormalizedPathPattern(path string, pattern string) (*regexp.Regexp, error) {
	preparedPath := NormalizeDirSep(path)
	preparedPatternToCompile := regexp.QuoteMeta(preparedPath) + "/" + pattern
	return regexp.Compile(preparedPatternToCompile)
}