package newdesignpattern

type ObserverInterface interface {
	Notify(args...interface{})
}