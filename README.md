# trpc

[![Build Status](https://travis-ci.org/shric/go-trpc.svg?branch=master)](https://travis-ci.org/shric/go-trpc)
[![codecov](https://codecov.io/gh/shric/go-trpc/branch/master/graph/badge.svg)](https://codecov.io/gh/shric/go-trpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/shric/go-trpc)](https://goreportcard.com/report/github.com/shric/go-trpc)
![MIT license](https://img.shields.io/github/license/shric/go-trpc)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shric/go-trpc)

A pleasant frontend for transmission.

This project was only started recently and therefore is very barebones.

The intention is to provide a more user friendly interface compared to `transmission-remote`.

## Features (unchecked will be implemented soon)

* Filters (applies to list and rm)
  * [x] -i, --incomplete: Include only incomplete torrents

* [ ] Pass filenames instead of IDs.

* list command
  * [x] Basic list functionality

* add command
  * [x] Add torrents by filename or URL
  * [x] Add paused torrents by URL
  * [ ] Add paused torrents by filename

* rm command
  * [x] --nuke: Remove local data as well as torrent.
  * [x] --force-all: Really remove all torrents if no IDs specified.


* [ ] start command

* [ ] stop command

## Installation

I wouldn't consider this very useful yet, but:

```sh
go install github.com/shric/go-trpc/cmd/trpc
```
## Usage

| trpc                    | transmission-remote analogue                      | Description                              |
| ----------------------- | ------------------------------------------------- | ---------------------------------------- |
| trpc -h                 | transmission-remote -h                            | Show help                                |
| trpc list               | transmission-remote -l                            | List all torrents                        |
| trpc list -i            |                                                   | List all incomplete torrents             |
| trpc rm 123 456         | transmission-remote -t 123,456 -r                 | Remove torrents with IDs 123 and 456     |
| trpc rm --nuke 123      | transmission-remote -t 123 -rad                   | Remove torrent + data with ID 123        |
| trpc add foo.torrent    | transmission-remote -a foo.torrent                | Add foo.torrent (can be filename or URL) |
| trpc add -p foo.torrent | transmission-remote -a --start-paused foo.torrent | Add foo.torrent in paused state          |
| trpc rm -i --force-all  |                                                   | Remove all incomplete torrents           |

