package file

import (
	"os"
	"bytes"
	"path/filepath"
	"github.com/sandreas/graft/pattern"
	"regexp"
	"errors"
	"io"
)
// Todo add progress callback parameter
func WalkPathByPattern(path string, compiledPattern *regexp.Regexp)([]string, error) {
	list := make([]string, 0)
	if path == "" {
		path = "."
	}
	err := filepath.Walk(path, func(innerPath string, info os.FileInfo, err error) error {
		normalizedPath := pattern.NormalizeDirSep(innerPath)
		// fmt.Println(" normalized ===> " + normalizedPath)
		if ! compiledPattern.MatchString(normalizedPath) {
			return nil
		}
		list = append(list, innerPath)
		return nil
	})

	//fmt.Println(list)
	return list, err
}

func Exists(f string) (bool) {
	_, err := os.Stat(f)
	if ! os.IsNotExist(err) {
		return true
	}
	return false
}

func ContentsEqual(src, dst string)(bool, error) {

	srcPointer, err := os.Open(src)
	if err != nil {
		return false, err
	}

	dstPointer, err := os.Open(dst)
	if err != nil {
		return false, err
	}

	return FileContentsEqual(srcPointer, dstPointer)
}

func FileContentsEqual(src, dst *os.File) (bool, error) {
	src.Seek(0, 0)
	dst.Seek(0, 0)
	srcBuffer := make([]byte, 32*1024)
	dstBuffer := make([]byte, 32*1024)
	for {
		n, err := src.Read(srcBuffer)
		if err != nil && err != io.EOF {
			return false, err
		}
		if n == 0 {
			o, _ := dst.Read(dstBuffer)
			// destination is larger than source
			if o != n {
				return false, nil
			}
			break
		}

		o, err := dst.Read(dstBuffer)
		if err != nil && err != io.EOF {
			return false, err
		}
		if o == 0 {
			break
		}

		if(!bytes.Equal(srcBuffer, dstBuffer)) {
			return false, nil
		}

	}

	return true, nil
}

func FileContentsEqualQuick(fi, fo *os.File, bufSize int64) (bool, error) {
	inStats, err := (fi).Stat()
	if(err != nil) {
		return false, err
	}
	outStats, err := fo.Stat()
	if(err != nil) {
		return false, err
	}


	inSize := inStats.Size()
	outSize := outStats.Size()

	if (outSize > inSize) {
		return false, nil
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

func Replace(src, dst string)(error) {
	_, err := os.Stat(dst)
	if ! os.IsNotExist(err) {
		os.Remove(dst)
	}
	return Copy(src, dst)
}


func Copy(src, dst string)(error) {
	srcStat, err := os.Stat(src)
	if os.IsNotExist(err) {
		return err
	}

	_, err = os.Stat(dst)

	if ! os.IsNotExist(err) {
		return errors.New("destination file " + dst + " already exists")
	}

	srcPointer, err := os.Open(src)
	if(err != nil) {
		return nil
	}

	dstPointer, _ := os.OpenFile(dst, os.O_RDWR | os.O_CREATE,srcStat.Mode())
	io.Copy(dstPointer, srcPointer)
	defer srcPointer.Close()
	defer dstPointer.Close()
	return nil
}

func CopyResumed(src, dst *os.File, progressHandler func(bytesTransferred, size, chunkSize int64) int64) (error) {
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

	bufferSize := progressHandler(dstSize, srcSize, 32*1024)
	buf := make([]byte, bufferSize)
	bytesTransferred := dstSize
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}


		if _, err := dst.Write(buf[:n]); err != nil {
			return err
		}
		bytesTransferred += int64(n);
		newBufferSize := progressHandler(bytesTransferred, srcSize, bufferSize)
		if(newBufferSize != bufferSize) {
			bufferSize = newBufferSize
			buf = make([]byte, bufferSize)
		}
	}
	return nil
}

//func MkdirAll(p string, perm os.FileMode) (error) {
//	pathParts := strings.Split(filepath.ToSlash(p), "/")
//	path := ""
//	for _, dir := range pathParts {
//		path += dir
//		stat, err := os.Stat(path)
//		if os.IsNotExist(err) {
//			os.Mkdir(path, perm)
//		} else if ! stat.IsDir() {
//			return errors.New("path " + path + " is a file and cannot be used as directory")
//		}
//
//		path += "/"
//
//	}
//	return nil
//}



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