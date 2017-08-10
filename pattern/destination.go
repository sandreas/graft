package pattern

import (
	"regexp"
)

type DestinationPattern struct {
	BasePattern
}

func NewDestinationPattern(patternString string) *DestinationPattern {
	destinationPattern := &DestinationPattern{}
	destinationPattern.parse(patternString)
	destinationPattern.fixRegex()
	return destinationPattern
}

// 	replace $1_ with ${1}_ to prevent problems during rename
func (p *DestinationPattern) fixRegex() {
	dollarUnderscore, _ := regexp.Compile("\\$([1-9][0-9]*)_")
	p.Pattern = dollarUnderscore.ReplaceAllString(p.Pattern, "${$1}_")
}
