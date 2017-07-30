package newprogress

type ProgressHandler interface {
	IncreaseItems()
	IncreaseMatches()
	Finish()
}

type WalkProgressHandler struct {
	itemCount      int64
	matchCount     int64
	forceShow      bool
	outputCallback func(format string, a ...interface{}) (int, error)
	Interval       int64
}

func NewWalkProgressHandler(handle func(format string, a ...interface{}) (int, error)) *WalkProgressHandler {
	return &WalkProgressHandler{
		Interval:       100,
		outputCallback: handle,
	}
}

func (ph *WalkProgressHandler) IncreaseItems() {
	ph.forceShow = ph.itemCount == 0 || ph.forceShow
	ph.itemCount++
	ph.showProgress()
}

func (ph *WalkProgressHandler) IncreaseMatches() {
	ph.forceShow = ph.matchCount == 0

	ph.matchCount++
	ph.IncreaseItems()
}

func (ph *WalkProgressHandler) showProgress() {
	if !ph.forceShow && ph.itemCount%ph.Interval != 0 {
		return
	}

	if ph.matchCount == 0 {
		ph.outputCallback("\rscanning - total: %d", ph.itemCount)
	} else {
		ph.outputCallback("\rscanning - total: %d,  matchCount: %d", ph.itemCount, ph.matchCount)
	}

	if ph.itemCount > ph.Interval*10 {
		ph.Interval = 500
	}
}

func (ph *WalkProgressHandler) Finish() {
	ph.forceShow = true
	ph.showProgress()
	ph.outputCallback("\n")
	//ph.outputCallback("")
}
