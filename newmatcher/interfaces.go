package newmatcher

type CompositeInterface interface{
	Add(child CompositeInterface)
}

type MatcherInterface interface {
	Matches(pattern interface{}) bool
}
