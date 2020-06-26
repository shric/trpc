// Package torrent provides a useful accessors for human readable output.
package torrent

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

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

// Status returns the status of a torrent.
func Status(torrent *transmissionrpc.Torrent) string {
	statusStrings := map[transmissionrpc.TorrentStatus]string{
		transmissionrpc.TorrentStatusStopped:      "Stopped",
		transmissionrpc.TorrentStatusCheckWait:    "To Hash",
		transmissionrpc.TorrentStatusCheck:        "Hashing",
		transmissionrpc.TorrentStatusDownloadWait: "Queued",
		transmissionrpc.TorrentStatusSeedWait:     "Queued",
		transmissionrpc.TorrentStatusIsolated:     "No peers",
	}

	if val, ok := statusStrings[*torrent.Status]; ok {
		return val
	}

	return ""
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

// Have returns the number of bytes downloaded so far.
func Have(t *transmissionrpc.Torrent) int64 {
	return int64(t.SizeWhenDone.Byte()) - *t.LeftUntilDone
}

// Progress returns the progress as a percentage of a download or hash check.
func Progress(t *transmissionrpc.Torrent) float64 {
	if *t.RecheckProgress != 0 {
		return 100.0 * *t.RecheckProgress
	}

	return 100.0 * float64(Have(t)) / t.SizeWhenDone.Byte()
}

// Ratio returns the upload/size ratio.
func Ratio(t *transmissionrpc.Torrent) float64 {
	return float64(*t.UploadedEver) / t.SizeWhenDone.Byte()
}

// Age returns the age of the torrent (the later of DoneDate and AddedDate).
func Age(t *transmissionrpc.Torrent) int64 {
	lastActivity := int64(math.Max(float64(t.DoneDate.Unix()), float64(t.AddedDate.Unix())))
	now := time.Now().Unix()

	return now - lastActivity
}

// Priority returns the priority of a torrent (low, medium, high).
func Priority(t *transmissionrpc.Torrent) string {
	priorities := []string{"low", "normal", "high"}
	// We add one because
	// https://github.com/transmission/transmission/blob/master/libtransmission/transmission.h
	//     TR_PRI_LOW = -1,
	//     TR_PRI_NORMAL = 0, /* since NORMAL is 0, memset initializes nicely */
	//     TR_PRI_HIGH = 1
	return priorities[*t.BandwidthPriority+1]
}

// TrackerShortName returns the configured short name of a torrent.
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
	torrent.Percent = Progress(result.original)
	torrent.Up = fmt.Sprintf("%7.1f", torrent.up/float64(KiB))
	torrent.Down = fmt.Sprintf("%7.1f", torrent.down/float64(KiB))
	torrent.UploadedEver += result.UploadedEver
	torrent.Ratio = Ratio(result.original)

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

	torrent.Percent = Progress(torrent.original)
	torrent.Eta = torrent.eta()
	torrent.up = float64(*torrent.original.RateUpload)
	torrent.Up = fmt.Sprintf("%7.1f", torrent.up/float64(KiB))
	torrent.down = float64(*torrent.original.RateDownload)
	torrent.Down = fmt.Sprintf("%7.1f", torrent.down/float64(KiB))
	torrent.Ratio = Ratio(torrent.original)
	torrent.Priority = Priority(torrent.original)
	torrent.Trackershortname = torrent.trackershortname(conf)

	status := Status(torrent.original)
	if status != "" {
		torrent.Up = status
		torrent.Down = status
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
