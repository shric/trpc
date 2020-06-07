# trpc

[![Build Status](https://travis-ci.org/shric/go-trpc.svg?branch=master)](https://travis-ci.org/shric/go-trpc)
[![codecov](https://codecov.io/gh/shric/go-trpc/branch/master/graph/badge.svg)](https://codecov.io/gh/shric/go-trpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/shric/go-trpc)](https://goreportcard.com/report/github.com/shric/go-trpc)
![MIT license](https://img.shields.io/github/license/shric/go-trpc)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shric/go-trpc)

A pleasant frontend for transmission.

This project was only started recently and therefore is very barebones.

## Features

* Filters (applies to list and rm)
  * [x] -i, --incomplete: Include only incomplete torrents

* list command
  * [x] Basic list functionality

* add command
  * [x] Add torrents by filename or URL
  * [x] Add paused torrents by URL
  * [ ] Add paused torrents by filename

* rm command
  * [x] --nuke: Remove local data as well as torrent.
  * [x] --force-all: Really remove all torrents if no IDs specified.

* [ ] Pass filenames instead of IDs.

* [ ] start command

* [ ] stop command
