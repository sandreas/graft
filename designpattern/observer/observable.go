package designpattern


type ObservableInterface interface {
	RegisterObserver(observerInterface ObserverInterface)
	notifyObservers(args...interface{})
}


type Observable struct {
	observers []ObserverInterface
	ObservableInterface
}

func (p *Observable) RegisterObserver(observer ObserverInterface) {
	p.observers = append(p.observers, observer)
}

func (p *Observable) NotifyObservers(args...interface{}) {
	for _, o := range p.observers {
		o.Notify(args...)
	}
}
