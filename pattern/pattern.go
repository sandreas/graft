package pattern

import (
	"os"
	"strings"
	"bytes"
	"regexp"
	"time"
	"strconv"
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

	if path == "" {
		path = "."
	}

	return path, pattern
}

func GlobToRegex(glob string) (string) {
	var buffer bytes.Buffer
	r := strings.NewReader(glob)

	escape := false
	//braceOpen := 0
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

		//if (r == '{') {
		//	buffer.WriteString("(")
		//	braceOpen++
		//	continue
		//}
		//
		//if (r == '}') {
		//	buffer.WriteString(")")
		//	braceOpen--
		//	continue
		//}
		//
		//if (r == ',' && braceOpen > 0) {
		//	buffer.WriteString("|")
		//	continue
		//}

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
	if preparedPath == "." {
		preparedPath = ""
	} else if preparedPath != "" {
		preparedPath = regexp.QuoteMeta(preparedPath) + "/"
	}
	preparedPatternToCompile :=  preparedPath + pattern
	return regexp.Compile(preparedPatternToCompile)
}

func StrToAge(t string, reference time.Time) (time.Time, error) {
	modifyPattern, err := regexp.Compile("^([+-]?[0-9]+)[\\s]*([a-zA-Z]+)$")
	if(err != nil) {
		return reference, err
	}

	if modifyPattern.MatchString(t) {
		submatches := modifyPattern.FindStringSubmatch(t)
		modifier, err := strconv.Atoi(submatches[1])
		if err != nil {
			return reference, err
		}

		// age must be negated
		modifier *= -1

		unit := strings.ToLower(submatches[2])

		if(strings.HasPrefix(unit, "d")) {
			return reference.AddDate(0, 0, modifier), nil
		}
		if(strings.HasPrefix(unit, "w")) {
			return reference.AddDate(0, 0, modifier * 7), nil
		}
		if(strings.HasPrefix(unit, "mon")) {
			return reference.AddDate(0, modifier, 0), nil
		}
		if(strings.HasPrefix(unit, "y")) {
			return reference.AddDate(modifier, 0, 0), nil
		}

		if strings.HasPrefix(unit, "ns") {
			unit = "ns"
		} else if strings.HasPrefix(unit, "us") || strings.HasPrefix(unit, "µs") {
			unit = "us"
		} else if strings.HasPrefix(unit, "ms") {
			unit = "ms"
		} else if strings.HasPrefix(unit, "s") {
			unit = "s"
		} else if strings.HasPrefix(unit, "m") {
			unit = "m"
		} else if strings.HasPrefix(unit, "h") {
			unit = "h"
		}

		d, err := time.ParseDuration(strconv.Itoa(modifier)+unit)
		if err != nil {
			return reference, err
		}

		return reference.Add(d), nil
	}

	fixedPattern, err := regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
	if fixedPattern.MatchString(t) {
		layout := "2006-01-02"
		return time.Parse(layout, t)
	}

	layout := "2006-01-02T15:04:05.000Z"
	return time.Parse(layout, t)
}