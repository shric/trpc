// Package torrent provides a useful accessors for human readable output.
package torrent

import (
	"fmt"

	"github.com/hekmon/cunits/v2"
	"github.com/hekmon/transmissionrpc"
)

type unitMap struct {
	Amount float64
	Name   string
}

func (torrent Torrent) eta() string {
	if *torrent.original.LeftUntilDone == 0 {
		return "Done"
	}
	if *torrent.original.Eta == -1 {
		return "Unknown"
	}
	units := []unitMap{
		{86400 * 365.25, "years"},
		{86400 * 365.25 / 12, "months"},
		{86400 * 7, "weeks"},
		{86400, "days"},
		{3600, "hours"},
		{60, "mins"},
	}
	eta := float64(*torrent.original.Eta)
	for _, u := range units {
		if eta > u.Amount*3 {
			return fmt.Sprintf("%d %s", int(eta/u.Amount), u.Name)
		}
	}
	return fmt.Sprintf("%d secs", int64(eta))
}

func (torrent Torrent) have() int64 {
	return int64(torrent.original.SizeWhenDone.Byte()) - *torrent.original.LeftUntilDone
}

func (torrent Torrent) progress() float64 {
	if torrent.original.SizeWhenDone == nil {
		return 100.0
	}
	if *torrent.original.RecheckProgress != 0 {
		return 100.0 * *torrent.original.RecheckProgress
	}
	return 100.0 * float64(torrent.have()) / float64(torrent.original.SizeWhenDone.Byte())
}

// NewFrom takes a transmissionrpc Torrent and provides useful human readable fields
func NewFrom(transmissionrpcTorrent *transmissionrpc.Torrent) (torrent *Torrent) {
	torrent = &Torrent{}
	torrent.original = transmissionrpcTorrent
	torrent.ID = *torrent.original.ID
	torrent.Name = *torrent.original.Name
	torrent.Size = *torrent.original.SizeWhenDone
	torrent.Error = " "
	if *torrent.original.Error != 0 {
		torrent.Error = "*"
	}
	torrent.Percent = torrent.progress()
	torrent.Eta = torrent.eta()
	return
}

// Torrent contains all the fields of transmissionrpc.Torrent but with non-pointer values
// useful for formatted output.
type Torrent struct {
	ID       int64
	Error    string
	Name     string
	Percent  float64
	Size     cunits.Bits
	Eta      string
	original *transmissionrpc.Torrent
}
