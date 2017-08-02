package newfile

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
	handler := NewWalkObserver(FakePrintf)
	handler.Interval = 2
	handler.Notify(LOCATOR_INCREASE_ITEMS)

	expect.Equal("\rscanning - total: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 1)

	handler.Notify(LOCATOR_INCREASE_ITEMS)

	expect.Equal("\rscanning - total: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 1)

	handler.Notify(LOCATOR_INCREASE_ITEMS)
	handler.Notify(LOCATOR_INCREASE_MATCHES)

	expect.Equal("\rscanning - total: %d,  matches: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 2)

	for i := 0; i < 20; i++ {
		handler.Notify(LOCATOR_INCREASE_ITEMS)
	}

	expect.Equal(int64(500), handler.Interval)

	handler.Notify(LOCATOR_FINISH)
	expect.Equal("\n", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 0)

}