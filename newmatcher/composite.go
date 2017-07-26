package newmatcher

type CompositeMatcher struct {
	CompositeInterface
	MatcherInterface
	children []MatcherInterface
}

func NewCompositeMatcher() *CompositeMatcher {
	return &CompositeMatcher{}
}

func (f *CompositeMatcher) Add(child MatcherInterface) {
	f.children = append(f.children, child)
}

func (f *CompositeMatcher) Matches(pattern string) bool {
	for _,val := range f.children {
		if ! val.Matches(pattern) {
			return false
		}
	}
	return true
}