package newtransfer

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

const DEFAULT_CHUNK_SIZE = int64(1024*32)
const DEFAULT_REPORT_INTERVAL = 300 * time.Millisecond

func TestEmptyFile(t *testing.T) {
	expect := assert.New(t)

	progressHandler := NewCopyProgressHandler(DEFAULT_CHUNK_SIZE, DEFAULT_REPORT_INTERVAL)
	newChunkSize, message := progressHandler.Update(0, 0, DEFAULT_CHUNK_SIZE, time.Now())
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("\r[====================>] 100.00%\n", message)
}

func TestNonEmptyFile(t *testing.T) {
	expect := assert.New(t)

	progressHandler := NewCopyProgressHandler(DEFAULT_CHUNK_SIZE, DEFAULT_REPORT_INTERVAL)

	size := int64(1024 * 1024 * 5)

	layout := "2006-01-02T15:04:05.000Z"
	nowAsString := "2017-08-02T21:45:00.000Z"
	now, _ := time.Parse(layout, nowAsString)

	newChunkSize, message := progressHandler.Update(0, size, DEFAULT_CHUNK_SIZE, now)
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("\r[>                    ] 0.00%", message)

	nowAsString = "2017-08-02T21:45:00.333Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered := int64(3*1024*1024)
	newChunkSize, message = progressHandler.Update(transfered, size, DEFAULT_CHUNK_SIZE, now)
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("\r[============>        ] 60.00%   9.01MiB/s", message)

	nowAsString = "2017-08-02T21:45:00.334Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered = int64(3*1024*1024 + 50)
	newChunkSize, message = progressHandler.Update(transfered, size, DEFAULT_CHUNK_SIZE, now)
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("", message)

	nowAsString = "2017-08-02T21:45:01.334Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered = int64(4*1024*1024 )
	newChunkSize, message = progressHandler.Update(transfered, size, DEFAULT_CHUNK_SIZE, now)
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("\r[================>    ] 80.00%   1022.98KiB/s", message)

	nowAsString = "2017-08-02T21:45:01.734Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered = size
	newChunkSize, message = progressHandler.Update(transfered, size, DEFAULT_CHUNK_SIZE, now)
	expect.Equal(DEFAULT_CHUNK_SIZE, newChunkSize)
	expect.Equal("\r[====================>] 100.00%   2.50MiB/s\n", message)
}


