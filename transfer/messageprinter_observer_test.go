package transfer_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/transfer"
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
	handler := transfer.NewMessagePrinterObserver(FakePrintf)

	handler.Notify("test-message")
	expect.Equal("test-message", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 0)

	handler.Notify("test-message with %s", "string")
	expect.Equal("test-message with %s", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 1)

	handler.Notify("test-message with a 100% percent sign and no parameters")
	expect.Equal("test-message with a 100%% percent sign and no parameters", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 0)
}