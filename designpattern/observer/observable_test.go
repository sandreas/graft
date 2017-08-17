package designpattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/designpattern/observer"
)

type FakeObserver struct {
	designpattern.ObservableInterface
	NotifyCalls int
	Argument int
}

func (fo *FakeObserver) Notify(a...interface{}) {
	fo.NotifyCalls++
	fo.Argument = a[0].(int)
}

func TestNewDestinationPattern(t *testing.T) {
	expect := assert.New(t)

	observable := &designpattern.Observable{}
	observer := &FakeObserver{}
	observable.RegisterObserver(observer)

	expect.Equal(0, observer.NotifyCalls)
	expect.Equal(0, observer.Argument)
	observable.NotifyObservers(1234)
	expect.Equal(1, observer.NotifyCalls)
	expect.Equal(1234, observer.Argument)
}