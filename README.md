# graft
graft is a command line application to search directories and transfer files. It started as a learning project and it still is, so much of the code could be vastly improved, but for now it already is a useful tool. 

## Features
- Finding and transferring files via glob patterns (`data/\*.jp\*g`) 
- Finding and transferring files via real regular expressions (`data/.\*\.jpe?g`)
- Provide additional filters like --max-age=2d (files older than 2 days are skipped)
- Copy and resume partially transferred files
- Exporting and importing file lists
- Providing files via sftp server

## Setup
**graft** should support Windows, MacOS and Linux, although the usage instructions might be different for each operating system. The easiest way to setup graft is to use the go package manager. See [installing go](https://golang.org/doc/install).

### First install
After [installing go](https://golang.org/doc/install) and adding the go binary to your Path, you can  now install graft with following command:

```
go get github.com/sandreas/graft
```

### Update

To update graft, simply use the `-u` flag
```
go get -u github.com/sandreas/graft
```

## Quickstart

###Important notes: 
- Every action is performed recursively by default, so all subdirectories are scanned
- It usually is a good idea to use the `--dry-run` option, to see, what graft is going to do with your files
- **Linux and Unix:** Use single quotes (') to encapsulate patterns 
- **Windows:** Use double quotes (") to encapsulate patterns

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
# move all jpeg files in /tmp/ to <filename>_new.<jpeg> (dry-run), e.g. /tmp/DSC0008.jpeg => /tmp/DSC0008_new.jpeg
graft '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' --move --dry-run
```

```
# copy all jpeg files in /tmp/ to  <filename>_new.<jpeg>, showing matching subexpressions
graft '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' --show-matches
```

## Usage Details

**graft** internally uses a combination of glob pattern conversion and regular expressions for matching and replacing file names.

### Basic Usage

Basic usage is:

```
graft [options] source [destination]
```

It usually is a good idea, to use the `--dry-run` option, to see, what graft is going to do with your files before something goes wrong.

### Notes for Windows and Unix 
Unix-Shells expand * and $1 by default, so use **single quotes**  (`'`) for all patterns to prevent unexpected results:

```
graft '/tmp/*.jpg' #
```

On Windows single quotes are treated as char, so use ***double quotes*** (`"`) and ***slashes*** (`/`) as directory separator:

```
graft "C:/Temp/*.jpg"
```

### Mode ***find***

The destination pattern is optional. If you do not specify a destination pattern, **graft** recursively lists all matching files and directories in all subdirectories, so it can also be used as a file search tool like find on unix systems. If you need to export these finds, there is an option for this.

#### Examples

Recursive listing of all jpg files in /tmp directory using a simple glob pattern:
```
graft '/tmp/*.jpg'
```

Using some regex-magic to find jpeg files, too:
```
graft '/tmp/*.jp[e]?g'
```

Exporting all results to a text file, one line for each find:
```
graft '/tmp/*.jpg' --export-to="~/jpg-in-tmp.txt"
```

Finding all files that are between 3 and 5 days old:
```
graft '/tmp/*.jpg' --min-age=3d --max-age=5d
```

You can also delete files - be careful with that... graft takes no prisoners and offers no apologies.
```
graft '/tmp/*.jpg' --min-age=3d --max-age=5d --delete
```

See **[Option reference](#option-reference)** for more info.

### Modes ***copy*** and ***move***

**graft** copies recursively and resumes partially transferred files by default. If you would like to move / rename files instead, use the ---move option 


Recursive copy every jpg file from tmp to /home/johndoe/pictures (dry-run)

```
graft '/tmp/*.jpg' '/home/johndoe/pictures/$1' --dry-run
```

#### Submatches and more complex examples 

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

If you would like to match ***()*** in directory names or file names, they have to be escaped via backslash:
```
graft '/tmp/videos \(2016\)' '/home/johndoe/'
```

You can also use pipes to match groups of chars:
```
graft '/tmp/(*.)(jpg|png)' '/home/johndoe/pictures/$1$2'
```

This will copy following source files to their destination:
```
/tmp/test.jpg          => /home/johndoe/pictures/test_new.jpg
/tmp/subdir/other.PNG  => /home/johndoe/pictures/subdir/other_new.PNG
```


### Option reference

Following options are available:
```
Positional arguments:
  SOURCE
  DESTINATION

Options:
  --move                 rename files instead of copy
  --delete               delete found files (be careful with this one - use --dry-run before execution)
  --dry-run              simulation mode - shows output but files remain unaffected
  --times                transfer source modify times to destination
  --force                force the requested action - even if it might be not a good idea
  --debug, -d            debug mode with logging to Stdout and into $HOME/.graft/application.log
  --quiet                do not show any output
  --show-matches         show pattern matches for each found file
  --export-to EXPORT-TO
                         export found matches to a text file - one line per item
  --files-from FILES-FROM
                         import found matches from file - one line per item
  --max-age MAX-AGE      maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )
  --min-age MIN-AGE      minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )
  --max-size MAX-SIZE    maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )
  --min-size MIN-SIZE    minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )
  --regex                use a real regex instead of glob patterns (e.g. src/.*\.jpg)
  --case-sensitive       be case sensitive when matching files and folders
  --sftpd                start sftp server providing only matching files and directories
  --sftp-password SFTP-PASSWORD
                         Specify the password for the sftp server
  --sftp-username SFTP-USERNAME
                         Specify the username for the sftp server [default: graft]
  --sftp-port SFTP-PORT
                         Specifies the port on which the server listens for connections [default: 2022]
  --help, -h             display this help and exit
```

The parameters `--min-age` and `--max-age` take duration or date strings to specify the age. Valid formats for age parameters, used like --min-age=X are:

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

The parameters `--min-size` and `--max-size` take size in bytes or size strings. Valid formats for size parameters, used like --min-size=X are:

```
1   => 1 byte
2M  => 2 MiB
3G  => 3 GiB
4T  => 4 TiB
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

**graft** is developed with JetBrains IntelliJ IDEA using the golang plugin, so this is the recommended IDE