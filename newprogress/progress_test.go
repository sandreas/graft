package newprogress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var lastFakePrintfString string
var lastFakePrintfParams []interface{}

func TestParse(t *testing.T) {
	expect := assert.New(t)
	handler := NewWalkProgressHandler(FakePrintf)
	handler.Interval = 2
	handler.IncreaseItems()

	expect.Equal("\rscanning - total: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 1)

	handler.IncreaseItems()

	expect.Equal("\rscanning - total: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 1)

	//handler.IncreaseItems()
	handler.IncreaseMatches()
	expect.Equal("\rscanning - total: %d,  matchCount: %d", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 2)

	for i := 0; i < 20; i++ {
		handler.IncreaseItems()
	}

	expect.Equal(int64(500), handler.Interval)

	handler.Finish()
	expect.Equal("\n", lastFakePrintfString)
	expect.Len(lastFakePrintfParams, 0)

}

func FakePrintf(format string, a ...interface{}) (int, error) {
	lastFakePrintfString = format
	lastFakePrintfParams = []interface{}{}
	for value := range a {
		lastFakePrintfParams = append(lastFakePrintfParams, value)
	}

	return 0, nil

}
