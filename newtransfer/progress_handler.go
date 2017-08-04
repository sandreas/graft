package newtransfer

import (
	"time"
	"github.com/tears-of-noobs/bytefmt"
	"math"
	"strconv"
	"fmt"

	"strings"
)

type CopyProgressHandler struct {
	//TransferProgressHandlerInterface
	bufferSize int64

	timerLastUpdate time.Time
	reportInterval time.Duration
	bytesLastUpdate int64

}

func NewCopyProgressHandler(bufSize int64, repInterval time.Duration) (*CopyProgressHandler) {
	return &CopyProgressHandler{
		bufferSize: bufSize,
		reportInterval: repInterval,
	}
}


func (s *CopyProgressHandler) Update(bytesTransferred, size, chunkSize int64, now time.Time) (int64, string) {
	s.startTimer(bytesTransferred, 1 * int64(time.Second), now)
	shouldReport := true
	bytesPerSecond := float64(0)
	percent := float64(0)

	shouldReport, bytesPerSecond, percent = s.getReportStatus(bytesTransferred, size, now)

	fmt.Printf("bytesPerSecond: %+v", bytesPerSecond)
	messageSuffix := ""

	if bytesTransferred == 0 {
		shouldReport = true
	}

	if bytesTransferred == size {
		shouldReport = true
		percent = 1
		messageSuffix = "\n"
	}

	if shouldReport {
		bandwidthOutput := "   " + bytefmt.FormatBytes(bytesPerSecond, 2, true) + "/s"
		charCountWhenFullyTransmitted := 20
		progressChars := int(math.Floor(percent * float64(charCountWhenFullyTransmitted)))
		normalizedInt := percent * 100
		percentOutput := strconv.FormatFloat(normalizedInt, 'f', 2, 64)
		if bytesPerSecond == 0 {
			bandwidthOutput = ""
		}
		//log.Printf("bandwidthOutput: %v, progressChars: %v, percentOutput: %v, messageSuffix: %v", bandwidthOutput, progressChars, percentOutput, messageSuffix)
		//return chunkSize, ""
		progressBar := fmt.Sprintf("[%-" + strconv.Itoa(charCountWhenFullyTransmitted + 1) + "s] " + percentOutput + "%%" + bandwidthOutput, strings.Repeat("=", progressChars) + ">")

		return chunkSize, "\r" + progressBar + messageSuffix
	}


	return chunkSize, ""
}


func (s *CopyProgressHandler) startTimer(bytesTransferred, interval int64, now time.Time) {
	if bytesTransferred == 0 {
		s.timerLastUpdate = now
		s.bytesLastUpdate = bytesTransferred
	}
}

func (s *CopyProgressHandler) getReportStatus(bytesTransferred, size int64, now time.Time) (bool, float64, float64) {
	nowNano := now.UnixNano()
	lastUpdateNano := s.timerLastUpdate.UnixNano()
	timeDiffNano := nowNano - lastUpdateNano
	timeDiffSeconds := float64(timeDiffNano) / float64(time.Second)
	// if timeDiffSeconds >= float64(reportInterval) {
	if timeDiffNano >= s.reportInterval.Nanoseconds() {
		bytesDiff := bytesTransferred - s.bytesLastUpdate
		bytesPerSecond := float64(float64(bytesDiff) / float64(timeDiffSeconds))
		percent := float64(bytesTransferred) / float64(size)

		s.bytesLastUpdate = bytesTransferred
		s.timerLastUpdate = now

		return true, bytesPerSecond, percent
	}

	return false, float64(0), float64(0)
}
