package torrent

import (
	"github.com/hekmon/cunits/v2"
	"github.com/hekmon/transmissionrpc"
)

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

	return torrent
}

// Torrent contains all the fields of transmissionrpc.Torrent but with non-pointer values
// useful for formatted output.
type Torrent struct {
	// These fields have been overridden so far.
	ID       int64
	Error    string
	Name     string
	Percent  float64
	Size     cunits.Bits
	original *transmissionrpc.Torrent
}
