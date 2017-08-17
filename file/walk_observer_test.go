package file_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/file"
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
	handler := file.NewWalkObserver(FakePrintf)
	handler.Interval = 2
	handler.Notify(file.LOCATOR_INCREASE_ITEMS)

	expect.Equal("\rscanning - total: %d,  matches: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 2)

	handler.Notify(file.LOCATOR_INCREASE_ITEMS)

	expect.Equal("\rscanning - total: %d,  matches: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 2)

	handler.Notify(file.LOCATOR_INCREASE_ITEMS)
	handler.Notify(file.LOCATOR_INCREASE_MATCHES)

	expect.Equal("\rscanning - total: %d,  matches: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 2)

	for i := 0; i < 20; i++ {
		handler.Notify(file.LOCATOR_INCREASE_ITEMS)
	}

	expect.Equal(int64(500), handler.Interval)

	handler.Notify(file.LOCATOR_FINISH)
	expect.Equal("\n", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 0)

}