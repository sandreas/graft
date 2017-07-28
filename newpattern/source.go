package newpattern

import (
	"os"
	"path/filepath"
	"strings"
	"regexp"
)

type Flag byte

const (
	CASE_SENSITIVE Flag = 1 << iota
	USE_REAL_REGEX
)

type SourcePattern struct {
	Path string
	Pattern string
	isDirectory bool
	caseSensitive bool
	useRealRegex bool
}

func NewSourcePattern(patternString string, params ...Flag) *SourcePattern {
	sourcePattern := &SourcePattern{}
	sourcePattern.parse(patternString)

	size := len(params)

	var flags Flag
	flags = 0x00
	for i := 0; i < size; i++ {
		flags |= params[i]
	}

	sourcePattern.caseSensitive = flags & CASE_SENSITIVE != 0
	sourcePattern.useRealRegex = flags & USE_REAL_REGEX != 0

	return sourcePattern
}


func (p *SourcePattern) parse(patternString string) {
	path := patternString
	for {
		if fi, err := os.Stat(path); err == nil {
			p.Path = filepath.ToSlash(path)

			p.isDirectory = fi.IsDir()


			startIndex := len(p.Path)+1
			if path == "." {
				p.Pattern = strings.TrimPrefix(patternString, ".")
				p.Pattern = strings.TrimPrefix(p.Pattern, "/")
			} else if len(patternString) > startIndex {
				p.Pattern = patternString[startIndex:]
			}
			break
		}
		path = filepath.Dir(path)
	}
}

func (p *SourcePattern) Compile() (*regexp.Regexp, error) {
	// pattern handling
	regexPattern := p.Pattern
	if ! p.useRealRegex {
		regexPattern = GlobToRegexString(p.Pattern)
	}
	if p.IsDir() && p.Pattern == "" {
		regexPattern = "(.*)"
	}

	// path handling
	regexPath := p.Path
	if regexPath != "" {
		regexPath = regexp.QuoteMeta(p.Path)
		if regexPath[len(regexPath)-1:] != "/" {
			regexPath += "/"
		}
	}

	if ! p.caseSensitive {
		regexPath = "(?i)" + regexPath
	}


	suffix := "$"
	compiledPattern, err := regexp.Compile(regexPath + regexPattern + suffix)
	if err == nil && compiledPattern.NumSubexp() == 0 && p.Pattern != "" {
		compiledPattern, err = regexp.Compile(regexPath + "(" + regexPattern + ")" + suffix)
	}

	return compiledPattern, err
}

func (p *SourcePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *SourcePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}