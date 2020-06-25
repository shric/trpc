// Package torrent provides a useful accessors for human readable output.
package torrent

import (
	"fmt"
	"net/url"
	"strconv"
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
	trackershortnameLen = 3
)

// Constants for binary units.
const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
)

func statusStr(torrent Torrent) *string {
	statusStrings := map[transmissionrpc.TorrentStatus]string{
		transmissionrpc.TorrentStatusStopped:      "Stopped",
		transmissionrpc.TorrentStatusCheckWait:    "To Hash",
		transmissionrpc.TorrentStatusCheck:        "Hashing",
		transmissionrpc.TorrentStatusDownloadWait: "Queued",
		transmissionrpc.TorrentStatusSeedWait:     "Queued",
		transmissionrpc.TorrentStatusIsolated:     "No peers",
	}

	if val, ok := statusStrings[*torrent.original.Status]; ok {
		return &val
	}

	return nil
}

func etastr(eta int64) string {
	units := []unitMap{
		{86400 * 365.25, "years"},
		{86400 * 365.25 / 12, "months"},
		{86400 * 7, "weeks"},
		{86400, "days"},
		{3600, "hours"},
		{60, "mins"},
	}

	secs := float64(eta)

	for _, u := range units {
		if secs > u.Amount*3 {
			return fmt.Sprintf("%d %s", int(secs/u.Amount), u.Name)
		}
	}

	return fmt.Sprintf("%d secs", int64(secs))
}

func (torrent Torrent) eta() string {
	if torrent.LeftUntilDone == 0 {
		return "Done"
	}

	if *torrent.original.Eta == -1 {
		return "Unknown"
	}

	return etastr(*torrent.original.Eta)
}

func (torrent Torrent) have() int64 {
	return int64(torrent.SizeWhenDone.Byte()) - torrent.LeftUntilDone
}

func (torrent Torrent) progress() float64 {
	if torrent.RecheckProgress != 0 {
		return 100.0 * torrent.RecheckProgress
	}

	return 100.0 * float64(torrent.have()) / torrent.SizeWhenDone.Byte()
}

func (torrent Torrent) ratio() float64 {
	// Returns +Inf on positive/0 or NaN on 0/0.
	return float64(torrent.UploadedEver) / torrent.SizeWhenDone.Byte()
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

func TrackerShortName(torrent *transmissionrpc.Torrent, conf *config.Config) string {
	for _, url := range torrent.Trackers {
		for match, shortname := range conf.Trackernames {
			if strings.Contains(url.Announce, match) {
				return shortname
			}
		}
	}
	return ""
}

func (torrent Torrent) trackershortname(conf *config.Config) string {
	if tsn := TrackerShortName(torrent.original, conf); tsn != "" {
		return tsn
	}

	url, err := url.Parse(torrent.original.Trackers[0].Announce)
	if err != nil {
		return "UNK"
	}

	hostname := url.Hostname()
	if len(hostname) >= trackershortnameLen {
		return url.Hostname()[0:trackershortnameLen]
	}

	return "UNK"
}

// UpdateTotal updates the Torrent carrying a total sum of torrents.
func (torrent *Torrent) UpdateTotal(result *Torrent) {
	torrent.up += result.up
	torrent.down += result.down
	torrent.SizeWhenDone += result.SizeWhenDone
	torrent.LeftUntilDone += result.LeftUntilDone
	torrent.Percent = torrent.progress()
	torrent.Up = fmt.Sprintf("%7.1f", torrent.up/float64(KiB))
	torrent.Down = fmt.Sprintf("%7.1f", torrent.down/float64(KiB))
	torrent.UploadedEver += result.UploadedEver
	torrent.Ratio = torrent.ratio()

	if torrent.LeftUntilDone != 0 && torrent.down != 0 {
		torrent.Eta = etastr(torrent.LeftUntilDone / int64(torrent.down))
	}
}

// NewForTotal returns a Torrent used to store the totals.
func NewForTotal() *Torrent {
	torrent := &Torrent{
		Error: " ",
	}

	return torrent
}

// NewFrom takes a transmissionrpc Torrent and provides useful human readable fields.
func NewFrom(transmissionrpcTorrent *transmissionrpc.Torrent, conf *config.Config) *Torrent {
	torrent := &Torrent{
		original:        transmissionrpcTorrent,
		ID:              strconv.FormatInt(*transmissionrpcTorrent.ID, 10),
		Name:            *transmissionrpcTorrent.Name,
		SizeWhenDone:    *transmissionrpcTorrent.SizeWhenDone,
		Status:          *transmissionrpcTorrent.Status,
		LeftUntilDone:   *transmissionrpcTorrent.LeftUntilDone,
		RecheckProgress: *transmissionrpcTorrent.RecheckProgress,
		UploadedEver:    *transmissionrpcTorrent.UploadedEver,
		Error:           " ",
	}

	if *torrent.original.Error != 0 {
		torrent.Error = "*"
	}

	torrent.Percent = torrent.progress()
	torrent.Eta = torrent.eta()
	torrent.up = float64(*torrent.original.RateUpload)
	torrent.Up = fmt.Sprintf("%7.1f", torrent.up/float64(KiB))
	torrent.down = float64(*torrent.original.RateDownload)
	torrent.Down = fmt.Sprintf("%7.1f", torrent.down/float64(KiB))
	torrent.Ratio = torrent.ratio()
	torrent.Priority = torrent.priority()
	torrent.Trackershortname = torrent.trackershortname(conf)

	status := statusStr(*torrent)
	if status != nil {
		torrent.Up = *status
		torrent.Down = *status
	}

	return torrent
}

// Torrent contains all the fields of transmissionrpc.Torrent but with non-pointer values
// useful for formatted output.
type Torrent struct {
	ID               string
	Error            string
	Name             string
	Percent          float64
	SizeWhenDone     cunits.Bits
	Eta              string
	up               float64
	Up               string
	down             float64
	Down             string
	Ratio            float64
	Priority         string
	Trackershortname string
	LeftUntilDone    int64
	RecheckProgress  float64
	UploadedEver     int64
	Status           transmissionrpc.TorrentStatus
	original         *transmissionrpc.Torrent
}
