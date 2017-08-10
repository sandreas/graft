- # graft
graft is a command line application to search directories and transfer files.
 
It supports glob patterns or regular expressions as well as resuming partial transferred files, preserving file attributes and exporting and importing lists.

***graft*** started as a learning project and it still is, so much of the code could be vastly improved and may contain bugs, 
but for now it already is a useful tool. 


## Installation

***graft*** should support Windows, MacOS and Linux, although the usage instructions might be different for each operating system. After [installing go](https://golang.org/doc/install) you can get graft with following command:

```
go get github.com/sandreas/graft
```

### Update

To update graft, simply use the *-u* flag
```
go get -u github.com/sandreas/graft
```

## Quickstart

###Important notes: 
- ***Linux and Unix:*** Use single quotes (') to encapsulate patterns 
- ***Windows:*** Use double quotes (") to encapsulate patterns

### Examples
```
# recursively search all jpg files in current directory and export a textfile
graft '*.jpg' --export-to=all-jpg-files.txt
```

```
# recursively copy all png files in data/ to /tmp
graft 'data/*.png' '/tmp/'
```

```
# start an sftp server promoting all txt files in data/ in a chroot 
graft 'data/*.txt' --serve --sftp-password=graft
```

```
# rename all jpeg files in data/ to jpg (dry-run)
graft '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' --dry-run

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

Recursive copy every jpg file from tmp to /home/johndoe/pictures (dry-run)

```
graft '/tmp/*.jpg' '/home/johndoe/pictures/$1' --dry-run
```

### Submatches and more complex examples 

As a result of using regular expressions, you can use `()` in combination with `$` to create submatches, e.g.:

```
graft '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2'
```

will copy following source files to their destination:

```
/tmp/test.jpeg          => /home/johndoe/pictures/test_new.jpeg
/tmp/subdir/other.jpeg  => /home/johndoe/pictures/subdir/other_new.jpeg
```
If you do not specify a submatch using ***()***, the whole pattern is treated as submatch.

```
graft '/tmp/*.jpg'

# is same as

graft '/tmp/(*.jpg)'
```

If you would like to match ***()*** in directory names or file names, they have to be escaped via backslash (\\):
```
graft '/tmp/videos \(2016\)' '/home/johndoe/'
```

You could also use brackets to match groups of chars:
```
graft '/tmp/(*.)(jpg|png)' '/home/johndoe/$1$2'
```

#### Option reference

Following options are available:
```
positional arguments:
  source
  destination

options:
  --case-sensitive       be case sensitive when matching files and folders
  --debug, -d            debug mode with logging to Stdout and into $HOME/.graft/application.log
  --delete               delete found files (be careful with this one - use --dry-run before execution)
  --dry-run              simulation mode output only files remain unaffected
  --move                 rename files instead of copy
  --regex                use a real regex instead of glob patterns (e.g. src/.*\.jpg)
  --quiet                do not show any output
  --sftp-promote         start sftp server only providing matching files and directories
  --show-matches         show pattern matches for each found file
  --times                transfer source modify times to destination
  --sftp-port SFTP-PORT
                         Specifies the port on which the server listens for connections (default: 2022) [default: 2022]
  --export-to EXPORT-TO
                         export found matches to a text file - one line per item
  --files-from FILES-FROM
                         import found matches from file - one line per item
  --max-age MAX-AGE      maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )
  --min-age MIN-AGE      minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )
  --sftp-password SFTP-PASSWORD
                         Specify the password for the sftp server
  --sftp-user SFTP-USER
                         Specify the username for the sftp server (default: graft) [default: graft]
  --help, -h             display this help and exit

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

***graft*** is developed with JetBrains IntelliJ IDEA using the golang plugin, so this is the recommended IDE