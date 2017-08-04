package newtransfer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var lastFakePrintfString string
var lastFakePrintfParams []interface{}

func FakePrintf(format string, a ...interface{}) (int, error) {
	lastFakePrintfString = format
	lastFakePrintfParams = []interface{}{}
	for value := range a {
		lastFakePrintfParams = append(lastFakePrintfParams, value)
	}

	return 0, nil

}

func TestParse(t *testing.T) {
	expect := assert.New(t)
	handler := NewMessagePrinterObserver(FakePrintf)
	handler.Notify("test-message")

	expect.Equal("test-message", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 0)

	handler.Notify("test-message with %s", "string")

	expect.Equal("test-message with %s", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 1)
}