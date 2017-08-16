package pattern

import (
	"regexp"
	"github.com/sandreas/graft/bitflag"
	"github.com/spf13/afero"
)


const (
	CASE_SENSITIVE bitflag.BitFlag = 1 << iota
	USE_REAL_REGEX
)

type SourcePattern struct {
	BasePattern
	caseSensitive bool
	useRealRegex bool
}

func NewSourcePattern(fs afero.Fs, patternString string, params ...bitflag.BitFlag) *SourcePattern {
	sourcePattern := &SourcePattern{}
	sourcePattern.Fs = fs
	sourcePattern.parse(patternString)


	bitFlags := bitflag.NewBitFlagParser(params...)
	sourcePattern.caseSensitive = bitFlags.HasFlag(CASE_SENSITIVE)
	sourcePattern.useRealRegex = bitFlags.HasFlag(USE_REAL_REGEX)

	return sourcePattern
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

	if regexPath == "." {
		regexPath = ""
	}

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
