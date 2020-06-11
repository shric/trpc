// Package torrent provides a useful accessors for human readable output.
package torrent

import (
	"fmt"
	"strings"

	"github.com/shric/trpc/internal/config"

	"github.com/hekmon/cunits/v2"
	"github.com/hekmon/transmissionrpc"
)

type unitMap struct {
	Amount float64
	Name   string
}

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
)

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
	if *torrent.original.RecheckProgress != 0 {
		return 100.0 * *torrent.original.RecheckProgress
	}

	return 100.0 * float64(torrent.have()) / torrent.original.SizeWhenDone.Byte()
}

func (torrent Torrent) ratio() float64 {
	// Returns +Inf on positive/0 or NaN on 0/0.
	return float64(*torrent.original.UploadedEver) / torrent.original.SizeWhenDone.Byte()
}

func (torrent Torrent) priority() string {
	priorities := []string{"low", "normal", "high"}
	// We add one because
	// https://github.com/transmission/transmission/blob/master/libtransmission/transmission.h
	//     TR_PRI_LOW = -1,
	//     TR_PRI_NORMAL = 0, /* since NORMAL is 0, memset initializes nicely */
	//     TR_PRI_HIGH = 1
	return priorities[*torrent.original.BandwidthPriority+1]
}

func (torrent Torrent) trackershortname(conf *config.Config) string {
	for _, url := range torrent.original.Trackers {
		for match, shortname := range conf.Trackernames {
			if strings.Contains(url.Announce, match) {
				return shortname
			}
		}
	}

	return "UNK"
}

// NewFrom takes a transmissionrpc Torrent and provides useful human readable fields.
func NewFrom(transmissionrpcTorrent *transmissionrpc.Torrent, conf *config.Config) *Torrent {
	torrent := &Torrent{
		original:      transmissionrpcTorrent,
		ID:            *transmissionrpcTorrent.ID,
		Name:          *transmissionrpcTorrent.Name,
		Size:          *transmissionrpcTorrent.SizeWhenDone,
		Status:        *transmissionrpcTorrent.Status,
		LeftUntilDone: *transmissionrpcTorrent.LeftUntilDone,
		Error:         " ",
	}

	if *torrent.original.Error != 0 {
		torrent.Error = "*"
	}

	torrent.Percent = torrent.progress()
	torrent.Eta = torrent.eta()
	torrent.Up = float64(*torrent.original.RateUpload) / float64(KiB)
	torrent.Down = float64(*torrent.original.RateDownload) / float64(KiB)
	torrent.Ratio = torrent.ratio()
	torrent.Priority = torrent.priority()
	torrent.Trackershortname = torrent.trackershortname(conf)

	return torrent
}

// Torrent contains all the fields of transmissionrpc.Torrent but with non-pointer values
// useful for formatted output.
type Torrent struct {
	ID               int64
	Error            string
	Name             string
	Percent          float64
	Size             cunits.Bits
	Eta              string
	Up               float64
	Down             float64
	Ratio            float64
	Priority         string
	Trackershortname string
	LeftUntilDone    int64
	Status           transmissionrpc.TorrentStatus
	original         *transmissionrpc.Torrent
}
