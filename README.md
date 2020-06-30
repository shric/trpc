# trpc

[![Build Status](https://travis-ci.org/shric/trpc.svg?branch=master)](https://travis-ci.org/shric/trpc)
[![codecov](https://codecov.io/gh/shric/trpc/branch/master/graph/badge.svg)](https://codecov.io/gh/shric/trpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/shric/trpc)](https://goreportcard.com/report/github.com/shric/trpc)
![MIT license](https://img.shields.io/github/license/shric/trpc)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shric/trpc)

A pleasant frontend for transmission.

This project was only started recently and is very much an alpha.
Backwards incompatible changes might be made, but in general there should be
no bugs on the master branch. I will sometimes push a chain of commits to master
where it may be broken in the middle but should be fine at the final master
commit.

The intention is to provide a more user friendly interface compared to
`transmission-remote`.

Here are some examples of things you can do easily that would require at least
some scripting with transmission-remote alone:

```sh
# Stop all incomplete torrents
trpc stop -i

# Start all torrents that live in a specific download dir
trpc start ~/torrent/foo/*
```


## Features

### Commands

`add`: add torrents by file or URL (-p to add paused)

`errors`: show torrents that have errors

`set`: Set global/torrent upload/download rate limits and torrent priorities

`list`: list torrents

`move`: move torrents to another location

`rename`: Rename a torrent path or file

`rm`: remove torrents (--nuke to delete the data as well as the torrent)

`start`: start torrents (--now to jump queue)

`stop`: stop torrents

`verify`: verify (hash check) torrents

`version`: show version

`watch`: Show progress for incomplete active torrents, quit when done

`which`: show which torrents (and file ID) given file(s) belong to.

### Filters

Most commands above except for `add`, `version`, `rename` and `which` accept
filter arguments to limit the torrents acted upon (or displayed in the case of
`list`):

`-i, --incomplete`: Include only incomplete torrents
`-a, --active`: Include only active torrents (downloading or uploading)
`-t, --tracker`: Match on tracker short name
`-e, --error`: Match on a specific error string
`-d, --download-dir`: Match on a download directory.

The above are all shorthand for a more powerful filter language:

```sh
# equivalent of trpc list -i
trpc list -f 'incomplete'
trpc list -f '!complete'

# equivalent of trpc list -a
trpc list -f 'up > 0 || down > 0'

# equivalent of trpc list -t foo
trpc list -f 'tracker == "foo"'

# equivalent of trpc list -e 'unregistered torrent'
trpc list -f 'error == "unregistered torrent"'

# equivalent of trpc list -d '/home/chris/images'
trpc list -f 'downloadDir == "/home/chris/images"'

# multiple expressions can be defined:
# List incomplete torrents larger than 1 GiB
trpc list -i -f 'size > 1 GiB'
trpc list -f '!complete' -f 'size > 1 GiB'
trpc list -f 'incomplete && size > 1 GiB'
```

### Sorting

* by size, name, id, ratio, age, have (amount of bytes downloaded), upload, progress
```sh
# Sort by size descending
trpc list --sort size -r
```

### Torrents can be selected by ID or filename

Unlike `transmission-remote`, you can refer to a torrent by its filename.
This allows easy shell globbing. Example:

```sh
# Pause all incomplete torrents in ~/torrents/recent
$ trpc stop --incomplete ~/torrent/recent/*
```
## Planned upcoming features (near future)

### More commands


`files`: List files within torrents

`get`: set specified files to be downloaded

`noget`: set specified files to not be downloaded

`info`: Show detailed torrent info

`sessioninfo`: Show session information



## Planned features (possible, distant future)

* A TUI mode

* Support bittorrent clients other than transmission

## Installation

```sh
go install github.com/shric/trpc/cmd/trpc
```

## More usage examples

| trpc                    | transmission-remote analogue (if present)         | Description                                    |
| ----------------------- | ------------------------------------------------- | ---------------------------------------------- |
| trpc -h                 | transmission-remote -h                            | Show help                                      |
| trpc list               | transmission-remote -l                            | List all torrents                              |
| trpc list -i            |                                                   | List all incomplete torrents                   |
| trpc rm 123 456         | transmission-remote -t 123,456 -r                 | Remove torrents with IDs 123 and 456           |
| trpc rm --nuke 123      | transmission-remote -t 123 -rad                   | Remove torrent + data with ID 123              |
| trpc add foo.torrent    | transmission-remote -a foo.torrent                | Add foo.torrent (can be filename or URL)       |
| trpc add -p foo.torrent | transmission-remote -a --start-paused foo.torrent | Add foo.torrent in paused state                |
| trpc rm -i --force-all  |                                                   | Remove all incomplete torrents                 |
| trpc start 123          | transmission-remote -t 123 -s                     | Start torrent 123                              |
| trpc start --now 123    |                                                   | Start torrent 123 (bypass queue)               |
| trpc stop 123           | transmission-remote -t 123 -S                     | Stop torrent 123                               |
| trpc list *             |                                                   | List all running torrents in the current dir   |
| trpc which filename.iso |                                                   | Identify which torrent filename.iso belongs to |
| trpc set --down 50 123  | transmission-remote -t 123 -d 50                  | Set torrent 123 download limit to 50KB/sec     |
| trpc set --down 0 123   | transmission-remote -t 123 -D                     | Remove download limit from torrent 123         |
| trpc set --down 50 -s   | transmission-remote -d 50                         | Set global download limit to 50KB/sec          |
| trpc set -p high 123    | transmission-remote -t 50 -Bh                     | Set torrent 125's bandwidth priority to high   |
