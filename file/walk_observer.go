package file

import "github.com/sandreas/graft/designpattern/observer"

type WalkObserver struct {
	designpattern.ObserverInterface
	itemCount      int64
	matchCount     int64
	forceShow      bool
	outputCallback func(format string, a ...interface{}) (int, error)
	Interval       int64
}

func NewWalkObserver(handle func(format string, a ...interface{}) (int, error)) *WalkObserver {
	return &WalkObserver{
		Interval:       100,
		outputCallback: handle,
	}
}

func (ph *WalkObserver) Notify(a...interface{}) {
	if a[0] == LOCATOR_INCREASE_ITEMS {
		ph.forceShow = ph.itemCount == 0
		ph.itemCount++
	}

	if a[0] == LOCATOR_INCREASE_MATCHES {
		ph.forceShow = ph.matchCount == 0
		ph.itemCount++
		ph.matchCount++
	}

	if a[0] == LOCATOR_FINISH {
		ph.forceShow = true
	}

	ph.showProgress()

	if a[0] == LOCATOR_FINISH {
		ph.outputCallback("\n")
	}
}

func (ph *WalkObserver) showProgress() {
	if !ph.forceShow && ph.itemCount%ph.Interval != 0 {
		return
	}

	if ph.matchCount == 0 {
		ph.outputCallback("\rscanning - total: %d", ph.itemCount)
	} else {
		ph.outputCallback("\rscanning - total: %d,  matches: %d", ph.itemCount, ph.matchCount)
	}

	if ph.itemCount > ph.Interval*10 {
		ph.Interval = 500
	}
}
