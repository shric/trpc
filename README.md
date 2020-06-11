# trpc

[![Build Status](https://travis-ci.org/shric/trpc.svg?branch=master)](https://travis-ci.org/shric/trpc)
[![codecov](https://codecov.io/gh/shric/trpc/branch/master/graph/badge.svg)](https://codecov.io/gh/shric/trpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/shric/trpc)](https://goreportcard.com/report/github.com/shric/trpc)
![MIT license](https://img.shields.io/github/license/shric/trpc)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shric/trpc)

A pleasant frontend for transmission.

This project was only started recently and therefore is very barebones.

The intention is to provide a more user friendly interface compared to `transmission-remote`.

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

`list`: list torrents

`move`: move torrents to another location

`rm`: remove torrents (--nuke to delete the data as well as the torrent)

`start`: start torrents (--now to jump queue)

`stop`: stop torrents

`verify`: verify (hash check) torrents

`version`: show version

### Filters

All commands above except for `add` and `version` accept filter arguments to limit the torrents acted upon (or displayed in the case of `list`):

`-i, --incomplete`: Include only incomplete torrents

### Torrents can be selected by ID or filename

Unlike `transmission-remote`, you can refer to a torrent by its filename. This allows easy shell globbing. Example:

```sh
# Pause all incomplete torrents in ~/torrents/recent
$ trpc stop --incomplete ~/torrent/recent/*
```
## Planned upcoming features (near future)

### More commands

`limit`: Set global/torrent upload/download rate limits

`files`: List files within torrents

`get`: set specified files to be downloaded

`noget`: set specified files to not be downloaded

`info`: Show detailed torrent info

`sessioninfo`: Show session information

`rename`: Rename the torrent path (without moving it to another downloadDir)

`watch`: Show a progress bar for incomplete active torrents

`which`: Identify which torrent a file belongs to

### More filters

* `-t, --tracker`: Select only torrents using a particular tracker

* `-e, --error`: Select only torrents with a given error substring

### Sorting

* by size

* by name

* by id

## Planned features (possible, distant future)

* A TUI mode

* Support bittorrent clients other than transmission

## Installation

```sh
go install github.com/shric/trpc/cmd/trpc
```

## More usage examples

| trpc                    | transmission-remote analogue (if present)         | Description                                  |
| ----------------------- | ------------------------------------------------- | -------------------------------------------- |
| trpc -h                 | transmission-remote -h                            | Show help                                    |
| trpc list               | transmission-remote -l                            | List all torrents                            |
| trpc list -i            |                                                   | List all incomplete torrents                 |
| trpc rm 123 456         | transmission-remote -t 123,456 -r                 | Remove torrents with IDs 123 and 456         |
| trpc rm --nuke 123      | transmission-remote -t 123 -rad                   | Remove torrent + data with ID 123            |
| trpc add foo.torrent    | transmission-remote -a foo.torrent                | Add foo.torrent (can be filename or URL)     |
| trpc add -p foo.torrent | transmission-remote -a --start-paused foo.torrent | Add foo.torrent in paused state              |
| trpc rm -i --force-all  |                                                   | Remove all incomplete torrents               |
| trpc start 123          | transmission-remote -t 123 -s                     | Start torrent 123                            |
| trpc start --now 123    |                                                   | Start torrent 123 (bypass queue)             |
| trpc stop 123           | transmission-remote -t 123 -S                     | Stop torrent 123                             |
| trpc list *             |                                                   | List all running torrents in the current dir |
