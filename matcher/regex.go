package matcher

import (
	"regexp"
)

type RegexMatcher struct {
	MatcherInterface
	regex *regexp.Regexp
}

func NewRegexMatcher(regex *regexp.Regexp) *RegexMatcher {
	return &RegexMatcher{
		regex: regex,
	}
}

func (f *RegexMatcher) Matches(subject interface{}) bool {
	return f.regex.MatchString(subject.(string))
}