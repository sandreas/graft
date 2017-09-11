# graft
graft is a command line utility to search and transfer files. It started as a learning project and it still is, so much of the code could be vastly improved, but for now it already is a useful tool. 

## Features
- Finding and transferring files via glob like patterns (`data/*.jp*g`) 
- Finding and transferring files via real regular expressions (`data/\.*\.jpe?g`)
- Provide additional filters like --max-age=2d (files older than 2 days are skipped)
- Copy and resume partially transferred files
- Exporting and importing file lists
- Providing and receive files over network via sftp server

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
- Every action is performed recursively by default, so all subdirectories are concerned
- For file creation commands, it usually is a good idea to use the `--dry-run` option, to see what graft is going to do
- **Linux and Unix:** Use single quotes (') to encapsulate patterns to prevent shell expansion
- **Windows:** Use double quotes (") to encapsulate patterns, since single quotes are treated as chars

### Examples
```
# recursively search all jpg files in current directory and export a textfile
graft find '*.jpg' --export-to=all-jpg-files.txt
```

```
# recursively copy all png files in data/ to /tmp
graft copy 'data/*.png' '/tmp/'
```

```
# start an sftp server promoting all txt files in data/ in a chroot environment via mdns/zeroconf
graft serve 'data/*.txt' --password=graft
```
```
# move all jpeg files in /tmp/ to <filename>_new.<jpeg> (dry-run), e.g. /tmp/DSC0008.jpeg => /tmp/DSC0008_new.jpeg
graft move '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' --dry-run
```

```
# copy all jpeg files in /tmp/ to  <filename>_new.<jpeg>
graft copy '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' 
```

### Network transfer
**graft** is designed to easily transfer files from a host to another. To achieve this, you can use `graft serve` and `graft receive` together.

To transfer all jpg files in `/tmp` from host A to host B, all you have to do is the following:

Host A:
```
graft serve '/tmp/*.jpg'
```

Host B:
```
graft receive
```

## Usage Details

**graft** internally uses a combination of glob pattern conversion and regular expressions for matching and replacing file names.


### ***find***

The find command is used to find files. In this mode **graft** recursively lists all matching files and directories in all subdirectories, so it can also be used as a file search tool like `find` on unix systems.

#### Examples

Recursive listing of all jpg files in /tmp directory using a simple glob pattern:
```
graft find '/tmp/*.jpg'
```

Using some regex-magic to find jpeg files, too:
```
graft find '/tmp/*.jp[e]?g'
```

Exporting all results to a text file, one line for each find:
```
graft find '/tmp/*.jpg' --export-to="~/jpg-in-tmp.txt"
```

Finding all files that are between 3 and 5 days old:
```
graft find '/tmp/*.jpg' --min-age=3d --max-age=5d
```

### ***serve***

The serve command is used to provide files via sftp. Similar to find, all matching files are provided via sftp. You can now use a sftp client like `Filezilla` or `WinSCP` to download files from the serving host.

Additionally, graft tries to use mdns/zeroconf to announce the sftp server within the current network, so that `graft receive` finds the server automatically and downloads all provided files.

Provide all jpg files in /tmp:
```
graft serve '/tmp/*.jpg'
```

By default graft serve will provide all files in the current directory:
```
graft serve

# is same as

graft serve .
```

### ***receive***
**graft** can receive files from a graft server. In most cases, it should find its pairing host automatically with zeroconf, but you can also specify, from which host you would like to receive.

Lookup host via zeroconf and receive files to current directory:
```
graft receive
```

Specify host:
```
graft receive --host 192.168.1.111
```

### ***delete***

You can also delete files - be careful with that... graft takes no prisoners and offers no apologies. By default you have to confirm the deletion.
```
graft delete '/tmp/*.jpg' --min-age=3d --max-age=5d
```

See **[Option reference](#option-reference)** for more info.



### ***copy***
**graft** can copy files recursively and resumes partially transferred files by default. 

Recursive copy every jpg file from `tmp` to `/home/johndoe/pictures (dry-run)

```
graft '/tmp/*.jpg' '/home/johndoe/pictures/$1' --dry-run
```

#### Submatches and more complex examples 

As a result of using regular expressions internally, you can use `()` in combination with `$` to create submatches, e.g.:

```
graft '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2'
```

will copy following source files to their destination:

```
/tmp/test.jpeg          => /home/johndoe/pictures/test_new.jpeg
/tmp/subdir/other.jpeg  => /home/johndoe/pictures/subdir/other_new.jpeg
```
If you do not specify a submatch using `(), the whole pattern is treated as submatch.

```
graft '/tmp/*.jpg'

# is same as

graft '/tmp/(*.jpg)'
```

If you would like to match `()` in directory names or file names, they have to be quoted via backslash:
```
graft '/tmp/videos \(2016\)' '/home/johndoe/'
```

You can also use pipes to match multiple variants of char combinations:
```
graft '/tmp/(*.)(jpg|png)' '/home/johndoe/pictures/$1$2'
```

This will copy following source files to their destination:
```
/tmp/test.jpg          => /home/johndoe/pictures/test_new.jpg
/tmp/subdir/other.PNG  => /home/johndoe/pictures/subdir/other_new.PNG
```

### ***move***
**graft** can also move files recursively. It should work exactly like copy except of moving files to its destination, instead of copy



### Reference

```
Usage: graft <action> SOURCE [OPTIONS]
  or   graft <action> SOURCE DESTINATION [OPTIONS] 

Positional arguments:
  SOURCE                Source file or directory
  DESTINATION           Destination file or directory (only available on transfer actions)
  
COMMANDS:
     find, f        find files
     serve, s       serve files
     copy, c, cp    copy files from a source to a destination
     move, m, mv    move files from a source to a destination
     delete, d, rm  delete files recursively
     receive, r     receive files from a graft server
     help, h        Shows a list of commands or help for one command


GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

GLOBAL ACTION OPTIONS:
   --quiet, -q         do not show any output
   --force, -f         force the requested action - even if it might be not a good idea
   --debug             debug mode with logging to Stdout and into $HOME/.graft/application.log
   --regex             use a real regex instead of glob patterns (e.g. src/.*\.jpg)
   --case-sensitive    be case sensitive when matching files and folders
   --max-age value     maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )
   --min-age value     minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )
   --max-size value    maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )
   --min-size value    minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )
   --export-to value   export found matches to a text file - one line per item (can also be used as save cache for large scans)
   --files-from value  import found matches from file - one line per item (can also be used as load cache for large scans)

FIND OPTIONS:
   --host value        Specify the hostname for the server (client mode only)
   --port value        Specifiy server port (used for server- and client mode) (default: 2022)
   --username value    Specify server username (used in server- and client mode) (default: "graft")
   --password value    Specify server password (used for server- and client mode)
   --show-matches      do not show matches for search pattern ($1=filename)
   --client            client mode - act as sftp client and search files remotely instead of local search


SERVE OPTIONS:
   --host value        Specify the hostname for the server (client mode only)
   --port value        Specifiy server port (used for server- and client mode) (default: 2022)
   --username value    Specify server username (used in server- and client mode) (default: "graft")
   --password value    Specify server password (used for server- and client mode)
   --no-zeroconf       do not use mdns/zeroconf to publish multicast sftp server (graft receive will not work without parameters)

COPY OPTIONS:
   --times             transfer source modify times to destination
   --dry-run           simulation mode - shows output but files remain unaffected
   
MOVE OPTIONS:
   --times             transfer source modify times to destination
   --dry-run           simulation mode - shows output but files remain unaffected   

DELETE OPTIONS:
   --dry-run           simulation mode - shows output but files remain unaffected

RECEIVE OPTIONS:
   --dry-run           simulation mode - shows output but files remain unaffected
   --times             transfer source modify times to destination
   --host value        Specify the hostname for the server (client mode only)
   --port value        Specifiy server port (used for server- and client mode) (default: 2022)
   --username value    Specify server username (used in server- and client mode) (default: "graft")
   --password value    Specify server password (used for server- and client mode)
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

If the build is successful, the directory should contain a binary named `graft` or `graft.exe` on windows systems
