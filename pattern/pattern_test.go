package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/pattern"
	"regexp"
	"time"
)

func TestNormalizeDirSep(t *testing.T) {
	expect := assert.New(t)
	expect.Equal("/tmp/dir/subdir", pattern.NormalizeDirSep("/tmp\\dir\\subdir"))
}

func TestParsePathPattern(t *testing.T) {
	expect := assert.New(t)

	path, pat := pattern.ParsePathPattern("../data/fixtures/global/*")
	expect.Equal("../data/fixtures/global", path)
	expect.Equal("*", pat)

	path, pat = pattern.ParsePathPattern("../data/fixtures/non-existing/*.*")
	expect.Equal("../data/fixtures", path)
	expect.Equal("non-existing/*.*", pat)

	path, pat = pattern.ParsePathPattern("../data/fixtures/global/file.txt")
	expect.Equal("../data/fixtures/global/file.txt", path)
	expect.Equal("", pat)

}

func TestGlobToRegex(t *testing.T) {
	expect := assert.New(t)
	expect.Equal(".*\\.jpg", pattern.GlobToRegex("*.jpg"))
	expect.Equal("star-file-\\*\\.jpg", pattern.GlobToRegex("star-file-\\*.jpg"))
	expect.Equal("test\\.(jpg|png)", pattern.GlobToRegex("test.{jpg,png}"))

	expect.Equal("fixtures\\(\\..*)", pattern.GlobToRegex("fixtures\\(.*)"))
}

func TestBuildMatchList(t *testing.T) {
	expect := assert.New(t)
	compiled, _ := regexp.Compile("data/fixtures/global/(.*)(\\.txt)$")

	list := pattern.BuildMatchList(compiled, "data/fixtures/global/documents (2010)/document (2010).txt")


	expect.Equal(2, len(list))
	expect.Equal("documents (2010)/document (2010)", list[0])
	expect.Equal(".txt", list[1])
}

func TestCompileNormalizedPathPattern(t *testing.T) {
	expect := assert.New(t)
	compiled, _ := pattern.CompileNormalizedPathPattern("data\\fixtures/global", "(.*)")
	expect.Equal("data/fixtures/global/(.*)", compiled.String())
	expect.Regexp(compiled, "data/fixtures/global/test.txt")


	compiled, _ = pattern.CompileNormalizedPathPattern("", "(.*\\.jpg)")
	expect.Equal("(.*\\.jpg)", compiled.String())
	expect.Regexp(compiled, "data/fixtures/global/test.jpg")

}

func TestStrToAge(t *testing.T) {
	expect := assert.New(t)

	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	referenceTime, _ := time.Parse(layout, str)

	actualTime, _ := pattern.StrToAge(str, time.Now())
	expectedTime  := referenceTime
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	tstr := "2014-11-12"
	actualTime, _ = pattern.StrToAge(tstr, time.Now())
	expectedTime, _  = time.Parse("2006-01-02", tstr)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("10 s", referenceTime)
	expectedTime  = referenceTime.Add(time.Second * -time.Duration(10))
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("10m", referenceTime)
	expectedTime  = referenceTime.Add(time.Minute * -time.Duration(10))
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("10h", referenceTime)
	expectedTime  = referenceTime.Add(time.Hour * -time.Duration(10))
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("1 day", referenceTime)
	expectedTime  = referenceTime.AddDate(0, 0, -1)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("-5 days", referenceTime)
	expectedTime  = referenceTime.AddDate(0, 0, 5)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("1 week", referenceTime)
	expectedTime  = referenceTime.AddDate(0, 0, -7)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("3 months", referenceTime)
	expectedTime  = referenceTime.AddDate(0, -3, 0)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = pattern.StrToAge("2 years", referenceTime)
	expectedTime  = referenceTime.AddDate(-2, 0, 0)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())
}
