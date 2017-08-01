package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newdesignpattern/observer"
	"path/filepath"
	"os"
	"github.com/sandreas/graft/newmatcher"
)

const (
	OBSERVER_INCREASE_ITEMS = 1
	OBSERVER_INCREASE_MATCHES = 2
	OBSERVER_FINISH = 3
)

type Transfer struct {
	newdesignpattern.Observable
	src newpattern.SourcePattern
	SourceFiles []string
}


func NewTransfer(pattern newpattern.SourcePattern) *Transfer {
	return &Transfer{
		src: pattern,
	}
}


func (t *Transfer) Find(matcher *newmatcher.CompositeMatcher) {
	t.SourceFiles = []string{}
	if t.src.IsFile() {
		t.SourceFiles = append(t.SourceFiles, t.src.Path)

		t.NotifyObservers(OBSERVER_INCREASE_MATCHES)
		t.NotifyObservers(OBSERVER_FINISH)
		return
	}

	filepath.Walk(t.src.Path, func(innerPath string, info os.FileInfo, err error) error {
		if innerPath == "." || innerPath == ".." {
			return nil
		}


		normalizedInnerPath := filepath.ToSlash(innerPath)
		if info.IsDir() {
			normalizedInnerPath += "/"
		}

		if matcher.Matches(normalizedInnerPath) {
			t.SourceFiles = append(t.SourceFiles, normalizedInnerPath)
			t.NotifyObservers(OBSERVER_INCREASE_MATCHES)
		} else {
			t.NotifyObservers(OBSERVER_INCREASE_ITEMS)
		}

		return nil
	})

	t.NotifyObservers(OBSERVER_FINISH)
}


func (t *Transfer) CopyTo(dst *newpattern.DestinationPattern) {

	if dst.IsFile() {

	}

}


//func (p *DestinationPattern) BuildDestination(element, patternPath string ) (string) {
//	var dst string
//	if p.Pattern == "" {
//		dst = filepath.ToSlash(p.Path + element[len(patternPath) + 1:])
//	} else {
//		dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(element), pattern.NormalizeDirSep(destinationPattern))
//	}
//
//	return dst
//}

/*
dstPath, dstPatt := pattern.ParsePathPattern(destinationPattern)

// replace $1_ with ${1}_ to prevent problems
dollarUnderscore, _ := regexp.Compile("\\$([1-9][0-9]*)_")
destinationPattern = dollarUnderscore.ReplaceAllString(destinationPattern, "${$1}_")

var dst string
for _, element := range matchingPaths {
	if dstPatt == "" {
		dst = pattern.NormalizeDirSep(dstPath + element[len(patternPath) + 1:])
	} else {
		dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(element), pattern.NormalizeDirSep(destinationPattern))
	}
	transferElementHandler(element, dst)
}

func transferElementHandler() {

	prntln(src + " => " + dst)

	if *dryRun {
		return
	}

	srcStat, srcErr := os.Stat(src)

	if srcErr != nil {
		prntln("could not read source: ", srcErr)
		return
	}

	dstStat, _ := os.Stat(dst)
	dstExists := file.Exists(dst)
	if srcStat.IsDir() {
		if ! dstExists {
			if os.MkdirAll(dst, srcStat.Mode()) != nil {
				prntln("Could not create destination directory")
			}
			appendRemoveDir(dst)
			fixTimes(dst, srcStat)
			return
		}

		if dstStat.IsDir() {
			appendRemoveDir(dst)
			fixTimes(dst, srcStat)
			return
		}

		prntln("destination already exists as file, source is a directory")
		return
	}

	if dstExists && dstStat.IsDir() {
		prntln("destination already exists as directory, source is a file")
		return
	}

	srcDir := path.Dir(src)
	srcDirStat, _ := os.Stat(srcDir)

	dstDir := path.Dir(dst)
	if ! file.Exists(dstDir) {
		os.MkdirAll(dstDir, srcDirStat.Mode())
	}

	if *move {
		renameErr := os.Rename(src, dst)
		if renameErr == nil {
			appendRemoveDir(srcDir)
			fixTimes(dst, srcStat)
			return
		}
		prntln("Could not rename source")
		return
	}

	srcPointer, srcPointerErr := os.Open(src)
	if srcPointerErr != nil {
		prntln("Could not open source file")
		return
	}
	dstPointer, dstPointerErr := os.OpenFile(dst, os.O_WRONLY | os.O_CREATE, srcStat.Mode())

	if dstPointerErr != nil {
		prntln("Could not create destination file", dstPointerErr.Error())
		return
	}

	file.CopyResumed(srcPointer, dstPointer, handleProgress)
	fixTimes(dst, srcStat)
}
 */


//func (t *Transfer) moveTo(dst string) {
//
//}
//
//func (t *Transfer) remove(dst string) {
//
//}