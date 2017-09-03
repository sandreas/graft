package transfer_test

import (
	"testing"

	"os"
	"time"

	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/transfer"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/designpattern/observer"
)

func TestRelativeWildcardMapping(t *testing.T) {
	expect := assert.New(t)

	strategy := prepareStrategy("src/*", "dst")

	expect.Equal("dst", strategy.DestinationFor("src"))
	expect.Equal("dst/src-file.txt", strategy.DestinationFor("src/src-file.txt"))
	expect.Equal("dst/test-dir", strategy.DestinationFor("src/test-dir"))
	expect.Equal("dst/test-dir/test-dir-file.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))

}

func TestComplexRelativeMapping(t *testing.T) {
	expect := assert.New(t)

	// copy data/fixtures/global/textfile.txt ../out/
var strategy *transfer.AbstractStrategy
	strategy = prepareStrategy("src/test-dir/test-dir-file.txt", "../out/")
	expect.Equal("../out/test-dir-file.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))


	strategy = prepareStrategy("src/test-dir", "../out")
	expect.Equal("../out/test-dir/test-dir-file.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))

	strategy = prepareStrategy("src/test-dir/test-dir-file.txt", "dst/overwrite.txt")
	expect.Equal("dst/overwrite.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))

	strategy = prepareStrategy("src/test-dir/test-dir-file.txt", "dst")
	expect.Equal("dst/test-dir-file.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))

}

func prepareStrategy(src string, dst string) *transfer.AbstractStrategy {
	fs := prepareFileSystem()
	srcPattern := pattern.NewSourcePattern(fs, src)
	dstPattern := pattern.NewDestinationPattern(fs, dst)
	compiledSrcPattern, _ := srcPattern.Compile()

	return &transfer.AbstractStrategy{
		SourcePattern:         srcPattern,
		DestinationPattern:    dstPattern,
		CompiledSourcePattern: compiledSrcPattern,
		KeepTimes:             false,
	}

}

func prepareFileSystem() afero.Fs {
	appFS := afero.NewMemMapFs()
	appFS.Mkdir("src", 0644)
	afero.WriteFile(appFS, "src/src-file.txt", []byte(""), 0755)
	appFS.Mkdir("src/test-dir", 0644)
	appFS.Mkdir("dst", 0644)
	afero.WriteFile(appFS, "dst/overwrite.txt", []byte(""), 0755)
	afero.WriteFile(appFS, "src/test-dir/test-dir-file.txt", []byte(""), 0755)
	return appFS
}

func TestSingleTransferSourceNotExists(t *testing.T) {
	expect := assert.New(t)
	strategy := prepareStrategy("src/*", "dst")
	expect.Error(strategy.PerformSingleTransfer("non-existing-file.txt"))
}

func TestSingleTransferDir(t *testing.T) {
	expect := assert.New(t)
	strategy := prepareStrategy("src/*", "dst")
	expect.NoError(strategy.PerformSingleTransfer("src/test-dir"))
	stat, err := strategy.DestinationPattern.Fs.Stat("dst/test-dir")
	expect.True(stat.IsDir())
	expect.NoError(err)
	expect.Len(strategy.TransferredDirectories, 1)
}

func TestSingleTransferFileDirectoryCreation(t *testing.T) {
	expect := assert.New(t)
	strategy := prepareStrategy("src/*", "dst")
	expect.Nil(strategy.PerformSingleTransfer("src/test-dir/test-dir-file.txt"))
	stat, err := strategy.DestinationPattern.Fs.Stat("dst/test-dir")
	expect.True(stat.IsDir())
	expect.NoError(err)
	expect.Len(strategy.TransferredDirectories, 0)
}

func TestSingleTransferTimes(t *testing.T) {
	expect := assert.New(t)

	layout := "2006-01-02T15:04:05.000Z"
	timeAsStr := "2014-11-12T11:45:26.371Z"
	referenceTime, _ := time.Parse(layout, timeAsStr)

	strategy := prepareStrategy("src/*", "dst")
	strategy.KeepTimes = true
	strategy.SourcePattern.Fs.Chtimes("src/test-dir", referenceTime, referenceTime)

	expect.Nil(strategy.PerformSingleTransfer("src/test-dir/test-dir-file.txt"))

	stat, _ := strategy.DestinationPattern.Fs.Stat("dst/test-dir")
	expect.Equal(referenceTime, stat.ModTime())
}

func TestMultiTransfer(t *testing.T) {
	expect := assert.New(t)
	strategy := prepareStrategy("src/*", "dst/$1")

	toTransfer := []string{"src", "src/test-dir"}

	expect.NoError(strategy.Perform(toTransfer))
	stat, err := strategy.DestinationPattern.Fs.Stat("dst/test-dir")
	expect.True(stat.IsDir())
	expect.NoError(err)
	expect.Len(strategy.TransferredDirectories, 2)
}

func TestDryRun(t *testing.T) {
	expect := assert.New(t)
	strategy := prepareStrategy("src/*", "dst/$1")
	strategy.DryRun = true
	toTransfer := []string{"src", "src/test-dir"}

	expect.NoError(strategy.Perform(toTransfer))
	stat, err := strategy.DestinationPattern.Fs.Stat("dst/test-dir")
	expect.Nil(stat)
	expect.True(os.IsNotExist(err))
	expect.Len(strategy.TransferredDirectories, 0)
}



type FakeObserver struct {
	designpattern.ObserverInterface
	messages []string
}

func (ph *FakeObserver) Notify(a ...interface{}) {
	ph.messages = append(ph.messages, a[0].(string))
}

func prepareFilesystemTest(src, srcContent, dst, dstContent string) (*pattern.SourcePattern, *pattern.DestinationPattern) {
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, src, []byte(srcContent), 0644)

	if dstContent != "" {
		afero.WriteFile(appFS, dst, []byte(dstContent), 0644)
	}

	return pattern.NewSourcePattern(appFS, ""), pattern.NewDestinationPattern(appFS, "")
}

func TestCopyNewFile(t *testing.T) {
	expect := assert.New(t)

	srcFile := "test1-src.txt"
	srcContents := "this is a file without existing destination"
	destinationFile := "test1-dst.txt"

	sourcePattern, destinationPattern := prepareFilesystemTest(srcFile, srcContents, destinationFile, "")


	subject,_ := transfer.NewTransferStrategy(transfer.Copy, sourcePattern, destinationPattern)
	subject.ProgressHandler = transfer.NewCopyProgressHandler(2, 1*time.Nanosecond)
	observer := &FakeObserver{}
	subject.RegisterObserver(observer)

	srcStats, _ := subject.SourcePattern.Fs.Stat(srcFile)
	err := subject.PerformFileTransfer(srcFile, destinationFile, srcStats)
	expect.Equal(nil, err)

	dstContents, _ := afero.ReadFile(subject.SourcePattern.Fs, destinationFile)
	expect.Equal(srcContents, string(dstContents))

	expect.Len(observer.messages, 2)
	expect.Equal("\r[>                    ] 0.00%", observer.messages[0])
	expect.Equal("\r[====================>] 100.00%", observer.messages[1][0:32])
}

func TestCopyLargerSourceError(t *testing.T) {
	expect := assert.New(t)

	srcFile := "test-src.txt"
	srcContents := "this is a small src with larger dst"
	destinationFile := "test-dst.txt"
	dstContents := "this is a dst that is larger than its source and therefore cannot be copied"
	sourcePattern, destinationPattern := prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)

	subject,_ := transfer.NewTransferStrategy(transfer.Copy, sourcePattern, destinationPattern)


	srcStats, _ := subject.SourcePattern.Fs.Stat(srcFile)
	err := subject.PerformFileTransfer(srcFile, destinationFile, srcStats)
	expect.Error(err)
	contents, _ := afero.ReadFile(subject.SourcePattern.Fs, destinationFile)
	expect.Equal(dstContents, string(contents))
}

func TestCopyPartial(t *testing.T) {
	expect := assert.New(t)

	srcFile := "test-src.txt"
	srcContents := "this is the full content of a file with a partial existing destination"
	destinationFile := "test-dst.txt"
	dstContents := "this is the full content of a file with a partial"
	sourcePattern, destinationPattern := prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)

	subject,_ := transfer.NewTransferStrategy(transfer.Copy, sourcePattern, destinationPattern)

	srcStats, _ := subject.SourcePattern.Fs.Stat(srcFile)
	err := subject.PerformFileTransfer(srcFile, destinationFile, srcStats)
	expect.Equal(nil, err)
	contents, _ := afero.ReadFile(subject.SourcePattern.Fs, destinationFile)
	expect.Equal(srcContents, string(contents))
}

func TestCopyExistingCompleted(t *testing.T) {
	expect := assert.New(t)
	srcFile := "test-src.txt"
	srcContents := "this is a file where src and dst are fully equal"
	destinationFile := "test-dst.txt"
	dstContents := "this is a file where src and dst are fully equal"
	sourcePattern, destinationPattern := prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)

	subject,_ := transfer.NewTransferStrategy(transfer.Copy, sourcePattern, destinationPattern)

	srcStats, _ := subject.SourcePattern.Fs.Stat(srcFile)
	err := subject.PerformFileTransfer(srcFile, destinationFile, srcStats)
	expect.Equal(nil, err)
	contents, _ := afero.ReadFile(subject.SourcePattern.Fs, destinationFile)
	expect.Equal(srcContents, string(contents))
}

func TestCopyZeroBytesFile(t *testing.T) {
	expect := assert.New(t)

	srcFile := "test-src.txt"
	srcContents := ""
	destinationFile := "test-dst.txt"
	dstContents := ""
	sourcePattern, destinationPattern := prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)

	subject,_ := transfer.NewTransferStrategy(transfer.Copy, sourcePattern, destinationPattern)

	srcStats, _ := subject.SourcePattern.Fs.Stat(srcFile)
	err := subject.PerformFileTransfer(srcFile, destinationFile, srcStats)
	expect.Equal(nil, err)
	contents, _ := afero.ReadFile(subject.SourcePattern.Fs, destinationFile)
	expect.Equal(srcContents, string(contents))
	_, err = subject.SourcePattern.Fs.Stat(destinationFile)
	expect.Nil(err)
}
func TestCleanupIsAlwaysNil(t *testing.T) {
	expect := assert.New(t)
	srcFile := "test-src.txt"
	srcContents := ""
	destinationFile := "test-dst.txt"
	dstContents := ""
	sourcePattern, destinationPattern := prepareFilesystemTest(srcFile, srcContents, destinationFile, dstContents)

	subject, _ := transfer.NewTransferStrategy(transfer.Copy, sourcePattern, destinationPattern)

	expect.Nil(subject.Cleanup())
}
