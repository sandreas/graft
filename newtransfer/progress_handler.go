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

	TimerLastUpdate time.Time
	reportInterval int64
	bytesLastUpdate int64
}

func NewCopyProgressHandler(bufSize int64) (*CopyProgressHandler) {
	return &CopyProgressHandler{
		bufferSize: bufSize,
	}
}


func (s *CopyProgressHandler) Update(bytesTransferred, size, chunkSize int64) (int64, string) {
	s.startTimer(bytesTransferred, 1 * int64(time.Second))
	shouldReport, bytesPerSecond, percent := s.getReportStatus(bytesTransferred, size)

	messageSuffix := ""
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
		progressBar := fmt.Sprintf("[%-" + strconv.Itoa(charCountWhenFullyTransmitted + 1) + "s] " + percentOutput + "%%" + bandwidthOutput, strings.Repeat("=", progressChars) + ">")

		return chunkSize, "\r" + progressBar + messageSuffix
	}


	return chunkSize, ""
}


func (s *CopyProgressHandler) startTimer(bytesTransferred, interval int64) {
	if bytesTransferred == 0 {
		s.TimerLastUpdate = time.Now()
		s.reportInterval = interval
		s.bytesLastUpdate = bytesTransferred
	}
}

func (s *CopyProgressHandler) getReportStatus(bytesTransferred, size int64) (bool, float64, float64) {
	timeDiffNano := time.Now().UnixNano() - s.TimerLastUpdate.UnixNano()
	timeDiffSeconds := float64(timeDiffNano) / float64(time.Second)
	// if timeDiffSeconds >= float64(reportInterval) {
	if timeDiffNano >= s.reportInterval {
		bytesDiff := bytesTransferred - s.bytesLastUpdate
		bytesPerSecond := float64(float64(bytesDiff) / float64(timeDiffSeconds))
		percent := float64(bytesTransferred) / float64(size)

		s.bytesLastUpdate = bytesTransferred
		s.TimerLastUpdate = time.Now()

		return true, bytesPerSecond, percent
	}

	return false, 0, 0
}


/*


func (s *CopyProgressHandler) update(bytesTransferred, size, chunkSize int64) int64 {
	if size <= 0 {
		return chunkSize
	}

	startTimer(bytesTransferred, 1 * int64(time.Second))
	shouldReport, bytesPerSecond, percent := getReportStatus(bytesTransferred, size)
	if shouldReport {
		bandwidthOutput := "   " + bytefmt.FormatBytes(bytesPerSecond, 2, true) + "/s"
		charCountWhenFullyTransmitted := 20
		progressChars := int(math.Floor(percent * float64(charCountWhenFullyTransmitted)))
		normalizedInt := percent * 100
		percentOutput := strconv.FormatFloat(normalizedInt, 'f', 2, 64)
		if bytesPerSecond == 0 {
			bandwidthOutput = ""
		}
		progressBar := fmt.Sprintf("[%-" + strconv.Itoa(charCountWhenFullyTransmitted + 1) + "s] " + percentOutput + "%%" + bandwidthOutput, strings.Repeat("=", progressChars) + ">")

		prnt("\r" + progressBar)
	}

	if bytesTransferred == size {
		prntln("")
	}

	return chunkSize
}

func getReportStatus(bytesTransferred, size int64) (bool, float64, float64) {
	timeDiffNano := time.Now().UnixNano() - timerLastUpdate.UnixNano()
	timeDiffSeconds := float64(timeDiffNano) / float64(time.Second)
	// if timeDiffSeconds >= float64(reportInterval) {
	if timeDiffNano >= reportInterval {
		bytesDiff := bytesTransferred - bytesLastUpdate
		bytesPerSecond := float64(float64(bytesDiff) / float64(timeDiffSeconds))
		percent := float64(bytesTransferred) / float64(size)

		bytesLastUpdate = bytesTransferred
		timerLastUpdate = time.Now()

		return true, bytesPerSecond, percent
	}

	return false, 0, 0
}
*/