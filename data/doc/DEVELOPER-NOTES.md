# Things to do

- Switch to subcommands with https://github.com/urfave/cli
    graft find
    graft copy
    graft delete
    graft move
    graft serve
    graft receive
    
    global flags
        MaxAge        string `arg:"--max-age,help:maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
        MinAge        string `arg:"--min-age,help:minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
        MaxSize       string `arg:"--max-size,help:maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
        MinSize       string `arg:"--min-size,help:minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
        Regex         bool `arg:"help:use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"`
        CaseSensitive 
        ShowMatches bool `arg:"--show-matches,help:show pattern matches for each found file"`
        ExportTo  string `arg:"--export-to,help:export found matches to a text file - one line per item"`
        FilesFrom string `arg:"--files-from,help:import found matches from file - one line per item"`



- add transfer.move_strategy_test
- add action.transfer_test

find
copy
move
delete
serve
receive


- update Matchers to use existing FileInfo for faster matching / use matcher.setFileInfo in FileMatcherInterface
- calculate and show transfer speed
- graft le/d.img ../out/$1 does not work
- fix possible concurrency problem with pathMapper
- --verbose (ls -lah like output)
- --files-only / --directories-only
- javascript plugins? https://github.com/robertkrimen/otto
- --hide-progress (for working like find)
- copy strategy:  ResumeSkipDifferent=default, ResumeReplaceDifferent (ReplaceAll, ReplaceExisting, SkipExisting)
- compare-strategy: quick, hash, full
- improve progress-bar output (progress speed is not accurate enough)
- sftp-server:
	    filezilla takes long and produces 0 byte files
		filesystem watcher for sftp server (https://godoc.org/github.com/fsnotify/fsnotify)
	accept connections from specific ip: 		conn, e := listener.Accept() clientAddr := conn.RemoteAddr() if clientAddr
- sftp client
  - mdns / bonjour client https://github.com/hashicorp/mdns
- --max-depth parameter (?)
- limit-results when searching or moving
- Input / Colors: https://github.com/dixonwille/wlog


# Technology links

## big list of different libraries
https://github.com/avelino/awesome-go#command-line

## command line parser
cli-app: https://github.com/urfave/cli
further info: https://nathanleclaire.com/blog/2014/08/31/why-codegangstas-cli-package-is-the-bomb-and-you-should-use-it/
nice: https://github.com/alecthomas/kingpin
https://github.com/alexflint/go-arg

## File copy

Bytewise copy: 
http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file-using-golang
http://stackoverflow.com/questions/20602131/io-writeseeker-and-io-readseeker-from-byte-or-file
http://stackoverflow.com/questions/38631982/golang-file-seek-and-file-writeat-not-working-as-expected

## file times
http://stackoverflow.com/questions/20875336/how-can-i-get-a-files-ctime-atime-mtime-and-change-them-using-golang

## Globbing
https://www.reddit.com/r/golang/comments/41ulfq/glob_for_go_works_much_faster_than_regexp_on/

# Tutorials

## Organize project structure
https://talks.golang.org/2014/organizeio.slide#22

## Unit Testing
https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742#.mco6oq8iu

## Regex
https://github.com/StefanSchroeder/Golang-Regex-Tutorial/blob/master/01-chapter2.markdown

## State of go
https://talks.golang.org/2017/state-of-go.slide#21