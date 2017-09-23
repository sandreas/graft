# graft
graft is a command line utility to search and transfer files. It started as a learning project and it still is, so much of the code could be vastly improved, but for now it already is a useful tool. 

## Features
- Finding and transferring files via glob like patterns (`data/*.jp*g`) 
- Finding and transferring files via real regular expressions (`data/\.*\.jpe?g`)
- Provide additional filters like --max-age=2d (files older than 2 days are skipped)
- Copy and resume partially transferred files
- Exporting and importing file lists
- Providing and receive files over network via sftp server

## Download and Setup
**graft** should support Windows, MacOS and Linux, although the usage instructions might be different for each operating system.

### Binary releases

**graft** is released as a single binary for all major platforms:


#### [💾 Windows Download](https://github.com/sandreas/graft/releases/download/v0.2.0/graft_0.2.0_windows_64bit.zip)

#### [💾 MacOS Download](https://github.com/sandreas/graft/releases/download/v0.2.0/graft_0.2.0_macOS_64bit.tar.gz)

#### [💾 Linux Download](https://github.com/sandreas/graft/releases/download/v0.2.0/graft_0.2.0_linux_64bit.tar.gz)

#### [Other releases on the release page](https://github.com/sandreas/graft/releases)

<script src="readme.js"></script>

### go get graft
If you would like to use **graft** on an unsupported platform, you can try the go package manager. 
After [installing go](https://golang.org/doc/install) and adding the go binary to your PATH, install graft with following command:

```
go get github.com/sandreas/graft
```

If compilation succeeds, you can use `graft` from the command line.

### Update via go get

To force an update of the graft sources, simply add the `-u` flag
```
go get -u github.com/sandreas/graft
```

## Quickstart

### Important notes: 
- Every action is performed recursively by default, so all subdirectories are concerned in every action
- For file transfer commands, it usually is a good idea to use the `--dry-run` option, to see what **graft** is going to do
- Special chars `\ . + * ? ( ) | [ ] { } ^ $` have to be quoted with backslash in patterns (e.g `graft find '/tmp/video*\(2016\)'`)
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
# move all jpeg files in /tmp/ to <filename>_new.<jpeg> (dry-run), e.g. /tmp/DSC0008.jpeg => /tmp/DSC0008_new.jpeg
graft move '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' --dry-run
```

```
# copy all jpeg files in /tmp/ to  <filename>_new.<jpeg>
graft copy '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2' 
```

```
# start an sftp server promoting all txt files in data/ in a chroot environment via zeroconf
graft serve 'data/*.txt' --password=graft
```

```
# Receive all files from a graft server running on 192.168.0.150
graft receive --host=192.168.0.150 --password=graft
```

### Network transfer
**graft** is designed for easy transferring files from one host to another. To achieve this, you can use `graft serve` and `graft receive`.

To transfer all jpg files in `/tmp` from **host A** to **host B**, all you have to do is the following:

Host A:
```
graft serve '/tmp/*.jpg'
```

**graft** will prompt for a password, run an sftp server and promote it via zeroconf. 

The sftp server uses following defaults:
- Port: 2022
- Username: graft
- Listen address: 0.0.0.0


To receive all files, all you have to do is:

```
graft receive
```

in the destination directory on any other host within the same network and type in your password. Partially copied files will be resumed.


## Usage Details

**graft** internally uses a combination of glob pattern conversion and regular expressions for matching and replacing file names.


### ***find***

The find command is used to find files. In this mode **graft** recursively lists all matching files and directories in all subdirectories, so it can also be used as a search tool like `find` on unix systems.

#### Examples

Recursive listing of all jpg files in /tmp directory using a simple glob pattern:
```
graft find '/tmp/*.jpg'
```

Using some regex-logic to find jpeg files, too:
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

See **[Option reference](#option-reference)** for more info.

### ***serve***

The serve command is used to provide files via sftp. Similar to find, all matching files are provided via sftp. You can now use a sftp client like `Filezilla` or `WinSCP` to download files from the serving host.

Additionally, graft uses mdns/zeroconf by default to announce the sftp server within the current network, so that `graft receive` finds the graft server automatically and downloads all provided files.

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

To login with `FileZilla, you have to use the correct protocol and port, e.g.:

- Server: sftp://192.168.0.150
- Username: graft (if you did not change the defaults)
- Password: <your password>
- Port: 2022 (if you did not change the defaults)

See **[Option reference](#option-reference)** for more info.


### ***receive***
**graft** can receive files from a graft server. It should find its pairing host automatically via zeroconf, but you can also specify the host to receive from.

Lookup host via zeroconf and receive files to current directory:

```
graft receive
```

Specify host and port:

```
graft receive --host=192.168.1.111 --port=2023
```

You can also specify a receive pattern and a destination pattern, to receive only matching files.

Receive only jpg files and put them into /tmp/jpgs/ (directory structure is preserved):

```
graft receive '*.jpg' '/tmp/jpgs/'
```

See **[Option reference](#option-reference)** for more info.


### ***delete***

You can also delete files. Be careful with this command. **graft** takes no prisoners and offers no apologies. 
Because of this, by default you have to confirm the deletion process. With `--force` the confirmation is skipped.

```
graft delete '/tmp/*.jpg' --min-age=3d --max-age=5d
```


See **[Option reference](#option-reference)** for more info.

### ***copy***
**graft** is a powerful copy tool. It can copy files recursively and resumes partially transferred files by default. 

Recursive copy every jpg file from `tmp` to `/home/johndoe/pictures` (dry-run)

```
graft copy '/tmp/*.jpg' '/home/johndoe/pictures/$1' --dry-run
```

#### Submatches and more complex examples 

As a result of using regular expressions internally, you can use `()` in combination with `$` to create submatches, e.g.:

```
graft copy '/tmp/(*).(jpeg)' '/home/johndoe/pictures/$1_new.$2'
```

will copy following source files to their destination:

```
/tmp/test.jpeg          => /home/johndoe/pictures/test_new.jpeg
/tmp/subdir/other.jpeg  => /home/johndoe/pictures/subdir/other_new.jpeg
```

If you do not specify a submatch using `()`, the whole pattern is treated as submatch.

```
graft copy '/tmp/*.jpg' '/tmp/copy/'

# is same as

graft copy '/tmp/(*.jpg)' '/tmp/copy/'
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

If you would like to match on of these chars `\ . + * ? ( ) | [ ] { } ^ $` in patterns, they have to be quoted via backslash:
```
graft copy '/tmp/*\(2016\)' '/home/johndoe/'
```

This means that on windows, you have to escape backslashes, when using patterns:
```
graft copy "folder\*example*\\*.jpg" "otherfolder\$1"
```
Will copy all jpg files of every subdirectory matching `example` to `otherfolder`. In existing directory names, backslashes need not to be escaped.
Since **graft** also works with slashes on windows, it is recommended to use slashes to prevent unexpected behaviour.


### ***move***
**graft** can also move files recursively. It works exactly like copy except of moving files to its destination, instead of making a copy.

### Reference

```
Usage: graft <action> SOURCE [OPTIONS]
  or   graft <action> SOURCE DESTINATION [OPTIONS] 

Positional arguments:
  SOURCE                Source file, directory or pattern
  DESTINATION           Destination file, directory or pattern (only available on transfer actions)
  
COMMANDS:
     find, f        find files
     serve, s       serve files via sftp server
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
   --show-matches      show regex matches for search pattern ($1=filename)
   --client            client mode - act as sftp client and search files remotely instead of local search
   --host value        Specify the hostname for the server (client mode only)
   --port value        Specifiy server port (used for server- and client mode) (default: 2022)
   --username value    Specify server username (used in server- and client mode) (default: "graft")
   --password value    Specify server password (used for server- and client mode)

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

***graft*** is developed in go and you can use the default `go build` command to create a working binary:

```
git clone https://github.com/sandreas/graft.git
cd graft

go build
```

If the build is successful, the directory should contain a binary named `graft` or `graft.exe` on windows systems


If you prefer to do a full release for all supported plattforms, use goreleaser:

```
git clone https://github.com/sandreas/graft.git
cd graft
go get github.com/goreleaser/goreleaser

# for stable releases
goreleaser 

# for current branch releases
goreleaser --snapshot
```

Your release files are placed in `dist`.