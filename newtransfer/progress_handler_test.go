package newtransfer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const DEFAULT_CHUNK_SIZE = int64(1024*32)

func TestEmptyFile(t *testing.T) {
	expect := assert.New(t)

	progressHandler := NewCopyProgressHandler(DEFAULT_CHUNK_SIZE)
	newChunkSize, message := progressHandler.Update(0, 0, DEFAULT_CHUNK_SIZE)
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("\r[====================>] 100.00%\n", message)
}


