package newpattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"regexp"
)

func TestGlobToRegex(t *testing.T) {
	expect := assert.New(t)
	expect.Equal(".*\\.jpg", GlobToRegexString("*.jpg"))
	expect.Equal("star-file-\\*\\.jpg", GlobToRegexString("star-file-\\*.jpg"))
	expect.Equal("test\\.(jpg|png)", GlobToRegexString("test.(jpg|png)"))
	expect.Equal("test\\.{1,}", GlobToRegexString("test.{1,}"))
	expect.Equal("fixtures\\(\\..*)", GlobToRegexString("fixtures\\(.*)"))
}


func TestStrToAge(t *testing.T) {
	expect := assert.New(t)

	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	referenceTime, _ := time.Parse(layout, str)

	actualTime, _ := StrToAge(str, time.Now())
	expectedTime  := referenceTime
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	tstr := "2014-11-12"
	actualTime, _ = StrToAge(tstr, time.Now())
	expectedTime, _  = time.Parse("2006-01-02", tstr)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("10 s", referenceTime)
	expectedTime  = referenceTime.Add(time.Second * -time.Duration(10))
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("10m", referenceTime)
	expectedTime  = referenceTime.Add(time.Minute * -time.Duration(10))
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("10h", referenceTime)
	expectedTime  = referenceTime.Add(time.Hour * -time.Duration(10))
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("1 day", referenceTime)
	expectedTime  = referenceTime.AddDate(0, 0, -1)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("-5 days", referenceTime)
	expectedTime  = referenceTime.AddDate(0, 0, 5)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("1 week", referenceTime)
	expectedTime  = referenceTime.AddDate(0, 0, -7)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("3 months", referenceTime)
	expectedTime  = referenceTime.AddDate(0, -3, 0)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())

	actualTime, _ = StrToAge("2 years", referenceTime)
	expectedTime  = referenceTime.AddDate(-2, 0, 0)
	expect.Equal(expectedTime.UnixNano(), actualTime.UnixNano())
}

func TestBuildMatchList(t *testing.T) {
	expect := assert.New(t)
	compiled, _ := regexp.Compile("data/fixtures/global/(.*)(\\.txt)$")

	list := BuildMatchList(compiled, "data/fixtures/global/documents (2010)/document (2010).txt")


	expect.Equal(2, len(list))
	expect.Equal("documents (2010)/document (2010)", list[0])
	expect.Equal(".txt", list[1])
}