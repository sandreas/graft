package newpattern

import (
	"bytes"
	"strings"
)

func GlobToRegexString(glob string) (string) {
	var buffer bytes.Buffer
	r := strings.NewReader(glob)

	escape := false
	for {
		r, _, err := r.ReadRune()
		if err != nil {
			break
		}

		if escape {
			buffer.WriteRune(r)
			escape = false
			continue
		}

		if r == '\\' {
			buffer.WriteRune(r)
			escape = true
			continue
		}

		if r == '*' {
			buffer.WriteString(".*")
			continue
		}

		if r == '.' {
			buffer.WriteString("\\.")
			continue
		}

		buffer.WriteRune(r)
	}

	return buffer.String()
}
