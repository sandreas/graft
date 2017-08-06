package newtransfer

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/afero"
	"github.com/sandreas/graft/newdesignpattern/observer"
	"time"
)

type FakeObserver struct {
	newdesignpattern.ObserverInterface
	messages []string
}

func(ph *FakeObserver) Notify(a...interface{}){
	ph.messages = append(ph.messages, a[0].(string))
}


func prepareFilesystemTest(src, srcContent, dst, dstContent string) afero.Fs {
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, src, []byte(srcContent), 0644)

	if dstContent != "" {
		afero.WriteFile(appFS, dst, []byte(dstContent), 0644)
	}
	return appFS
}

func TestCopyNewFile(t *testing.T) {
	expect := assert.New(t)

	subject := NewCopyStrategy()
	subject.SetProgressHandler(NewCopyProgressHandler(2, 1 * time.Nanosecond))
	observer := &FakeObserver{}
	subject.RegisterObserver(observer)

	srcFile := "test1-src.txt"
	srcContents := "this is a file without existing destination"
	destinationFile := "test1-dst.txt"

	subject.Fs = prepareFilesystemTest(srcFile, srcContents, destinationFile, "")
	err := subject.Transfer(srcFile, destinationFile)
	expect.Equal(nil, err)

	dstContents, _ := afero.ReadFile(subject.Fs, destinationFile)
	expect.Equal(srcContents, string(dstContents))

	expect.Len(observer.messages, 1)
	expect.Equal("\r[====================>] 100.00%   0.00/s\n", observer.messages[0])
}

func TestCopyLargerSourceError(t *testing.T) {
	expect := assert.New(t)

	subject := NewCopyStrategy()

	srcFile := "test-src.txt"
	srcContents := "this is a small src with larger dst"
	destinationFile := "test-dst.txt"
	dstContents := "this is a dst that is larger than its source and therefore cannot be copied"
	subject.Fs = prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)
	err := subject.Transfer(srcFile, destinationFile)
	expect.Error(err)
	contents, _ := afero.ReadFile(subject.Fs, destinationFile)
	expect.Equal(dstContents, string(contents))
}

func TestCopyPartial(t *testing.T) {
	expect := assert.New(t)

	subject := NewCopyStrategy()

	srcFile := "test-src.txt"
	srcContents := "this is the full content of a file with a partial existing destination"
	destinationFile := "test-dst.txt"
	dstContents := "this is the full content of a file with a partial"
	subject.Fs = prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)
	err := subject.Transfer(srcFile, destinationFile)
	expect.Equal(nil, err)
	contents, _ := afero.ReadFile(subject.Fs, destinationFile)
	expect.Equal(srcContents, string(contents))
}

func TestCopyExistingCompleted(t *testing.T) {
	expect := assert.New(t)

	subject := NewCopyStrategy()

	srcFile := "test-src.txt"
	srcContents := "this is a file where src and dst are fully equal"
	destinationFile := "test-dst.txt"
	dstContents := "this is a file where src and dst are fully equal"
	subject.Fs = prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)
	err := subject.Transfer(srcFile, destinationFile)
	expect.Equal(nil, err)
	contents, _ := afero.ReadFile(subject.Fs, destinationFile)
	expect.Equal(srcContents, string(contents))
}