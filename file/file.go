package file

import (
	"os"
	"bytes"
	"path/filepath"
	"github.com/sandreas/graft/pattern"
	"regexp"
	"fmt"
	"io"
	"errors"
)

//import "os"

func WalkPathByPattern(path string, compiledPattern *regexp.Regexp)([]string, error) {
	list := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		normalizedPath := pattern.NormalizeDirSep(path)
		if ! compiledPattern.MatchString(normalizedPath) {
			return nil
		}
		list = append(list, path)
		return nil
	})
	return list, err
}
/*
func CopyResumed(src, dst *os.File, progressHandler func(bytesTransferred, size int64) int64) (error) {
	srcStats, err := (*src).Stat()
	if err != nil {
		return err
	}

	dstStats, err := (*dst).Stat()
	if err != nil {
		return err
	}

	srcSize := srcStats.Size()
	dstSize := dstStats.Size()

	if(dstSize > srcSize) {
		return errors.New("File cannot be resumed, destination is larger than source")
	}

	if(srcSize == dstSize) {
		return nil
	}
	src.Seek(dstSize, 0)
	dst.Seek(dstSize, 0)

	bufferSize := 1024
	buf := make([]byte, bufferSize)
	for {
		// read a chunk
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := dst.Write(buf[:n]); err != nil {
			return err
		}
		newBufferSize := progressHandler(srcSize)
		if(newBufferSize != bufferSize) {
			bufferSize = newBufferSize
			buf = make([]byte, bufferSize)
		}
	}
	return nil
}
*/

func FilesEqualQuick(inFile, outFile string, bufSize int64) (bool, error) {
	inStats, err := os.Stat(inFile)
	if(err != nil) {
		return false, err
	}
	outStats, err := os.Stat(outFile)
	if(err != nil) {
		return false, err
	}


	inSize := inStats.Size()
	outSize := outStats.Size()

	if (outSize > inSize) {
		return false, nil
	}

	fi, err := os.Open(inFile)
	if err != nil {
		return false, err
	}


	fo, err := os.Open(outFile)
	if err != nil {
		return false, err
	}
	backBufSize := bufSize
	if bufSize > outSize {
		bufSize = outSize
		backBufSize = 0
	} else if outSize < bufSize * 2 {
		backBufSize = outSize - bufSize
	}

	fiBuf := make([]byte, bufSize)
	_, err = fi.ReadAt(fiBuf, 0)

	if err != nil {
		return false, err
	}

	foBuf := make([]byte, bufSize)
	_, err = fo.ReadAt(foBuf, 0)

	if err != nil {
		return false, err
	}

	if ! bytes.Equal(fiBuf, foBuf) {
		return false, nil
	}

	if backBufSize > 0 {
		backOffset := outSize - backBufSize
		fiBuf = make([]byte, backBufSize)
		_, err = fi.ReadAt(fiBuf, backOffset)
		if err != nil {
			return false, err
		}
		foBuf = make([]byte, backBufSize)
		_, err = fo.ReadAt(foBuf, backOffset)
		if err != nil {
			return false, err
		}
		if ! bytes.Equal(fiBuf, foBuf) {
			return false, nil
		}
	}

	return true, nil
}


//func TransferFileBuffer(fi, fo *os.File, offset int64, progressHandler func(bytesTransferred, size int64) int64) (error) {
//	var chunkSize int64
//	chunkSize = 1024
//	buf := make([]byte, chunkSize)
//
//	fi.Seek(offset, 0)
//	for {
//		// read a chunk
//		n, err := fi.Read(buf)
//		if err != nil && err != io.EOF {
//			return err
//		}
//		if n == 0 {
//			break
//		}
//
//		// write a chunk
//		if _, err := fo.Write(buf[:n]); err != nil {
//			return err
//		}
//
//		//bytesTransferred := offset + bytesWritten
//		//newChunkSize := progressHandler()
//
//	}
//}
//func CopyResumed(src, dst string, checkEqual func(in, out *os.File) bool, progressHandler func(bytesTransferred, size, chunkSize int64) int64) {
//
//}

//func (re *Regexp) CopyResume(src string, repl func(string) string) string {
//	b := re.replaceAll(nil, src, 2, func(dst []byte, match []int) []byte {
//		return append(dst, repl(src[match[0]:match[1]])...)
//	})
//	return string(b)
//}