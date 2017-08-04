package newtransfer

import "github.com/sandreas/graft/newdesignpattern/observer"

type MessagePrinterObserver struct {
	newdesignpattern.ObserverInterface
	outputCallback func(format string, a ...interface{}) (int, error)
}

func NewMessagePrinterObserver(handle func(format string, a ...interface{}) (int, error)) *MessagePrinterObserver {
	return &MessagePrinterObserver{
		outputCallback: handle,
	}
}

func (ph *MessagePrinterObserver) Notify(a...interface{}) {
	var str string
	var params[]interface{}
	a_len := len(a)
	if a_len > 0 {
		str = a[0].(string)
	}
	if a_len > 1 {
		params = a[1:]
	}
	ph.outputCallback(str, params...)
}
