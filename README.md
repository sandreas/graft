# graft
graft is a command line application written in go to search directories and transfer files.
 
It supports glob patterns or regular expressions as well as resuming partial transferred files, preserving file attributes and exporting and importing lists.

***graft*** started as a learning project and it still is, so much of the code could be vastly improved and may contain bugs, 
but for now it already is a useful tool that works well in most cases. 


## Installation

***graft*** should support Windows, MacOS and Linux, although the usage instructions might be different for each operating system. You can download the latest pre-compiled binary from the [release page](releases/latest) or if you already installed go development tools, install graft via:

```
go get github.com/sandreas/graft
```

## Usage

***graft*** internally uses a combination of globbing conversion and regular expressions for matching and replacing directory patterns.

### Basic Usage

Basic usage of graft is:

```
graft [options] source [destination]
```

### Search mode

The destination pattern is optional, as well as the other programm options. If you do not specify a destination pattern, ***graft*** recursively lists all matching files in all subdirectories, so it can also be used as a search tool.

### Copy, Rename, Resume

***graft*** copies recursively and resumes partially transferred files by default. If you would like to move / rename files instead, use the ---move option 


### Notes for Windows vs. Unix 
Unix-Shells expand * and $1 by default, so use ***single quotes*** (') for all patterns to prevent unexpected results:

```
graft '/tmp/*.jpg'
```

On Windows, use ***double quotes*** (") and ***slashes*** (/) as directory separator:

```
graft "C:/Temp/*.jpg"
```

It usually is a good idea, to use the ***--dry-run*** option, to see, what graft is going to do with your files.


### Simple Examples

Recursive listing of all jpg files in /tmp directory using a simple glob pattern:
 
```
graft '/tmp/*.jpg'
```

Recursive copy every jpg file from tmp to /home/johndoe/pictures

```
graft '/tmp/*.jpeg' '/home/johndoe/pictures/$1.jpeg'
```

Recursive rename all files with extension jpeg to jpg:

```
graft '/tmp/*.jpeg' '/tmp/$1.jpg' --move
```

### Submatches and more complex examples 

As a result of using regular expressions, you can use `()` in combination with `$` to create submatches, e.g.:

```
graft '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2'
```

will copy following source files to their destination:

/tmp/test.jpeg          => /home/johndoe/pictures/test_new.jpeg
/tmp/subdir/other.jpeg  => /home/johndoe/pictures/subdir/other_new.jpeg


If you would like to match `()` in directorynames or filenames, they have to be escaped via backslash (\\):
```
graft '/tmp/videos \(2016\)' '/home/johndoe/'
```

You could also use braces to match groups of chars:
```
graft '/tmp/*.{jpg,png}' '/home/johndoe/$1'
```

#### Option reference

Following options are available:
```
Flags:
  --help            Show context-sensitive help (also try --help-long and --help-man).
  --export-to=""    export source listing to file, one line per item
  --files-from=""   import source listing from file, one line per item
  --min-age=""      minimum age (e.g. 2d, 8w, 2016-12-24, etc. - see docs for valid time formats)
  --max-age=""      maximum age (e.g. 2d, 8w, 2016-12-24, etc. - see docs for valid time formats)
  --case-sensitive  be case sensitive when matching files and folders
  --dry-run         dry-run / simulation mode
  --hide-matches    hide matches in search mode ($1: ...)
  --move            move / rename files - do not make a copy
  --quiet           quiet mode - do not show any output
  --regex           use a real regex instead of glob patterns (e.g. src/.*\.jpg)
  --times           transfer source modify times to destination

Args:
  <source-pattern>         source pattern - used to locate files (e.g. src/*)
  [<destination-pattern>]  destination pattern for transfer (e.g. dst/$1)

```

The parameters --min-age and --max-age take duration or date strings to specify the age. Valid formats for age parameters, used like --min-age=X are:

```
1s                          => 1 second
2m                          => 2 minutes
3h                          => 3 hours
4d                          => 4 days
5w                          => 5 weeks
6mon                        => 6 months
7y                          => 7 years
2006-01-02                  => exact date 2006-01-02 00:00:00
2006-01-02T15:04:05.000Z    => exact date 2006-01-02 15:04:05
```

# development

***graft*** is developed go and uses the default go build command

```
git clone https://github.com/sandreas/graft.git

cd graft

go build
```

If the build is successful, the directory should contain a binary named `graft`

## IDE recommendation

***graft*** is developed with JetBrains IntelliJ IDEA, so this is the recommended IDE