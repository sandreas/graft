package designpattern

type ObserverInterface interface {
	Notify(args...interface{})
}