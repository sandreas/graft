package transfer_test

import (
	"testing"

	"github.com/sandreas/graft/pattern"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/transfer"
	"time"
	"os"
)

func TestRelativeWildcardMapping(t *testing.T) {
	expect := assert.New(t)

	strategy := prepareStrategy("src/*", "dst")

	expect.Equal("dst", strategy.DestinationFor("src"))
	expect.Equal("dst/src-file.txt", strategy.DestinationFor("src/src-file.txt"))
	expect.Equal("dst/test-dir", strategy.DestinationFor("src/test-dir"))
	expect.Equal("dst/test-dir/test-dir-file.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))



}

func TestComplexRelativeWildcardMapping(t *testing.T) {
	expect := assert.New(t)

	// copy data/fixtures/global/textfile.txt ../out/

	strategy := prepareStrategy("src/test-dir/test-dir-file.txt", "../out/")

	expect.Equal("../out/test-dir-file.txt", strategy.DestinationFor("src/test-dir/test-dir-file.txt"))

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
		KeepTimes: false,
	}

}

func prepareFileSystem() afero.Fs {
	appFS := afero.NewMemMapFs()
	appFS.Mkdir("src", 0644)
	afero.WriteFile(appFS, "src/src-file.txt", []byte(""), 0755)
	appFS.Mkdir("src/test-dir", 0644)
	appFS.Mkdir("dst", 0644)
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
	expect.Error(strategy.PerformSingleTransfer("src/test-dir/test-dir-file.txt"), "method PerformFileTransfer is abstract and must be overridden in strategy")
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

	expect.Error(strategy.PerformSingleTransfer("src/test-dir/test-dir-file.txt"), "method PerformFileTransfer is abstract and must be overridden in strategy")

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

