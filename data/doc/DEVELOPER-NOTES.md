# Things to do

- switch mdns library to https://github.com/grandcat/zeroconf

- global
    - choose a release strategy (https://github.com/goreleaser/goreleaser)
    - update documentation
    - graft archive command
    - find a device id to ensure that "move" is safe to do / otherwise use copy and delete afterwards
    - replace pathmapper with file_tree
    - close source fs
    - add matcher for mimetype (image/*, image/jpeg)
    - max-depth parameter (?)
    - hide progress?!
    - improve performance of huge amounts of small files
    - shouldStop return parameter for filesystem.Walk
        - limit-results parameter

- copy
    - support copy strategy:  ResumeSkipDifferent=default, ResumeReplaceDifferent (ReplaceAll, ReplaceExisting, SkipExisting)
    - support file compare stitching (reading first, last and middle bytes)
    - compare-strategy: quick, hash, full

- serve
    - supportmultiple mdns entries - switch mdns library to https://github.com/grandcat/zeroconf
    - Improve handling of huge amounts of files
    
    
Possible improvements
- update Matchers to use existing FileInfo for faster matching / use matcher.setFileInfo in FileMatcherInterface
- calculate and show transfer speed
- --verbose (ls -lah like output)
- --files-only / --directories-only
- javascript plugins? https://github.com/robertkrimen/otto
- --hide-progress (for working like find)
- improve progress-bar output (progress speed is not accurate enough)
- sftp-server:
		filesystem watcher for sftp server (https://godoc.org/github.com/fsnotify/fsnotify)
	accept connections from specific ip: 		conn, e := listener.Accept() clientAddr := conn.RemoteAddr() if clientAddr
- sftp client
  - mdns / bonjour client https://github.com/hashicorp/mdns
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