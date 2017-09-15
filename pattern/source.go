package pattern

import (
	"path/filepath"
	"regexp"

	"os"

	"github.com/sandreas/graft/bitflag"
	"github.com/spf13/afero"
	"strings"
)

const (
	CASE_SENSITIVE bitflag.Flag = 1 << iota
	USE_REAL_REGEX
)

type SourcePattern struct {
	BasePattern
	caseSensitive bool
	useRealRegex  bool
}

func NewSourcePattern(fs afero.Fs, patternString string, params ...bitflag.Flag) *SourcePattern {
	sourcePattern := &SourcePattern{}
	sourcePattern.Fs = fs
	sourcePattern.parse(patternString)

	bitFlags := bitflag.NewParser(params...)
	sourcePattern.caseSensitive = bitFlags.HasFlag(CASE_SENSITIVE)
	sourcePattern.useRealRegex = bitFlags.HasFlag(USE_REAL_REGEX)

	return sourcePattern
}

func (p *SourcePattern) Compile() (*regexp.Regexp, error) {
	// pattern handling
	regexPattern := p.Pattern
	if !p.useRealRegex {
		regexPattern = GlobToRegexString(p.Pattern)
	}
	if p.IsDir() && p.Pattern == "" {
		regexPattern = "(.*)"
	}

	// path handling
	regexPath := strings.TrimPrefix(p.Path, "." + string(os.PathSeparator))

	if regexPath != "" {
		regexPath = filepath.ToSlash(p.Path)
		if regexPath[len(regexPath)-1] != '/' && !p.IsFile() {
			regexPath += "/"
		}
		regexPath = regexp.QuoteMeta(regexPath)
	}

	if !p.caseSensitive {
		regexPath = "(?i)" + regexPath
	}

	// replace double path separator with single slash
	r := regexp.MustCompile("[" +regexp.QuoteMeta(string(os.PathSeparator)) + "/]{2,}")
	regexPattern = r.ReplaceAllStringFunc(regexPattern, func(m string) string {
		return "/"
	})

	suffix := "$"
	compiledPattern, err := regexp.Compile(regexPath + regexPattern + suffix)
	if err == nil && compiledPattern.NumSubexp() == 0 && p.Pattern != "" {
		compiledPattern, err = regexp.Compile(regexPath + "(" + regexPattern + ")" + suffix)
	}

	return compiledPattern, err
}
