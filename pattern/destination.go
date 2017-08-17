package pattern

import (
	"regexp"
	"github.com/spf13/afero"
)

type DestinationPattern struct {
	BasePattern
}

func NewDestinationPattern(fs afero.Fs, patternString string) *DestinationPattern {
	destinationPattern := &DestinationPattern{}
	destinationPattern.Fs = fs
	destinationPattern.parse(patternString)
	destinationPattern.fixRegex()
	return destinationPattern
}

// 	replace $1_ with ${1}_ to prevent problems during rename
func (p *DestinationPattern) fixRegex() {
	dollarUnderscore, _ := regexp.Compile("\\$([1-9][0-9]*)_")
	p.Pattern = dollarUnderscore.ReplaceAllString(p.Pattern, "${$1}_")
}
