package file

import "github.com/sandreas/graft/designpattern/observer"

type WalkObserver struct {
	designpattern.ObserverInterface
	itemCount      int64
	matchCount     int64
	errorCount     int64
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

func (ph *WalkObserver) Notify(a ...interface{}) {
	if a[0] == LocatorIncreaseItems {
		ph.forceShow = ph.itemCount == 0
		ph.itemCount++
	}

	if a[0] == LocatorIncreaseErrors {
		ph.forceShow = ph.errorCount == 0
		ph.errorCount++
	}

	if a[0] == LocatorIncreaseMatches {
		ph.forceShow = ph.matchCount == 0
		ph.itemCount++
		ph.matchCount++
	}

	if a[0] == LocatorFinish {
		ph.forceShow = true
	}

	ph.showProgress()

	if a[0] == LocatorFinish {
		ph.outputCallback("\n")
	}
}

func (ph *WalkObserver) showProgress() {
	if !ph.forceShow && ph.itemCount%ph.Interval != 0 {
		return
	}

	if ph.errorCount == 0 {
		ph.outputCallback("\rscanning - total: %d,  matches: %d", ph.itemCount, ph.matchCount)
	} else {
		ph.outputCallback("\rscanning - total: %d,  matches: %d, errors: %d", ph.itemCount, ph.matchCount, ph.errorCount)
	}

	if ph.itemCount > ph.Interval*10 {
		ph.Interval = 500
	}
}
