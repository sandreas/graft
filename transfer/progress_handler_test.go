package transfer_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/sandreas/graft/transfer"
)

const (
	defaultChunkSize      = int64(1024 * 32)
	defaultReportInterval = 300 * time.Millisecond
)

func TestEmptyFile(t *testing.T) {
	expect := assert.New(t)

	progressHandler := transfer.NewCopyProgressHandler(defaultChunkSize, defaultReportInterval)
	newChunkSize, message := progressHandler.Update(0, 0, defaultChunkSize, time.Now())
	expect.Equal(defaultChunkSize, newChunkSize)
	expect.Equal("\r[====================>] 100.00%\n", message)
}

func TestNonEmptyFile(t *testing.T) {
	expect := assert.New(t)

	progressHandler := transfer.NewCopyProgressHandler(defaultChunkSize, defaultReportInterval)

	size := int64(1024 * 1024 * 5)

	layout := "2006-01-02T15:04:05.000Z"
	nowAsString := "2017-08-02T21:45:00.000Z"
	now, _ := time.Parse(layout, nowAsString)

	newChunkSize, message := progressHandler.Update(0, size, defaultChunkSize, now)
	expect.Equal(defaultChunkSize, newChunkSize)
	expect.Equal("\r[>                    ] 0.00%", message)

	nowAsString = "2017-08-02T21:45:00.333Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered := int64(3 * 1024 * 1024)
	newChunkSize, message = progressHandler.Update(transfered, size, defaultChunkSize, now)
	expect.Equal(defaultChunkSize, newChunkSize)
	expect.Equal("\r[============>        ] 60.00%   9.01MiB/s", message)

	nowAsString = "2017-08-02T21:45:00.334Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered = int64(3*1024*1024 + 50)
	newChunkSize, message = progressHandler.Update(transfered, size, defaultChunkSize, now)
	expect.Equal(defaultChunkSize, newChunkSize)
	expect.Equal("", message)

	nowAsString = "2017-08-02T21:45:01.334Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered = int64(4 * 1024 * 1024)
	newChunkSize, message = progressHandler.Update(transfered, size, defaultChunkSize, now)
	expect.Equal(defaultChunkSize, newChunkSize)
	expect.Equal("\r[================>    ] 80.00%   1022.98KiB/s", message)

	nowAsString = "2017-08-02T21:45:01.734Z"
	now, _ = time.Parse(layout, nowAsString)
	transfered = size
	newChunkSize, message = progressHandler.Update(transfered, size, defaultChunkSize, now)
	expect.Equal(defaultChunkSize, newChunkSize)
	expect.Equal("\r[====================>] 100.00%   2.50MiB/s\n", message)
}
