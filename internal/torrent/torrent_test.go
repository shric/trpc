package torrent_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/shric/trpc/internal/torrent"

	"github.com/hekmon/cunits/v2"
	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/config"
)

func makeTorrent(sizeWhenDone cunits.Bits, eta int64, recheckProgress float64, leftUntilDone int64,
	uploadedEver int64, tracker string, error int64,
) *torrent.Torrent {
	ID := int64(1)
	Name := "Torrent 1"
	IsFinished := false
	RateUpload := int64(0)
	RateDownload := int64(0)
	BandwidthPriority := int64(0)
	Status := transmissionrpc.TorrentStatus(0)
	Trackers := []*transmissionrpc.Tracker{
		{Announce: tracker},
	}
	trpcTorrent := transmissionrpc.Torrent{
		SizeWhenDone:      &sizeWhenDone,
		ID:                &ID,
		Name:              &Name,
		Error:             &error,
		Eta:               &eta,
		IsFinished:        &IsFinished,
		RateUpload:        &RateUpload,
		RateDownload:      &RateDownload,
		RecheckProgress:   &recheckProgress,
		LeftUntilDone:     &leftUntilDone,
		Status:            &Status,
		UploadedEver:      &uploadedEver,
		BandwidthPriority: &BandwidthPriority,
		Trackers:          Trackers,
	}
	conf := config.Config{
		Trackernames: map[string]string{
			"foo-tracker": "foo",
		},
	}

	return torrent.NewFrom(&trpcTorrent, &conf)
}

type test struct {
	input *torrent.Torrent
	want  interface{}
}

func TestNewFromEta(t *testing.T) {
	tests := []test{
		{input: makeTorrent(cunits.Bits(1), 1, 0.0, 1, 0, "", 0), want: "1 secs"},
		{input: makeTorrent(cunits.Bits(1), 180, 0.0, 1, 0, "", 0), want: "180 secs"},
		{input: makeTorrent(cunits.Bits(1), 240, 0.0, 1, 0, "", 0), want: "4 mins"},
		{input: makeTorrent(cunits.Bits(1), 10000, 0.0, 1, 0, "", 0), want: "166 mins"},
		{input: makeTorrent(cunits.Bits(1), 20000, 0.0, 1, 0, "", 0), want: "5 hours"},
		{input: makeTorrent(cunits.Bits(1), 100000, 0.0, 1, 0, "", 0), want: "27 hours"},
		{input: makeTorrent(cunits.Bits(1), 1000000, 0.0, 1, 0, "", 0), want: "11 days"},
		{input: makeTorrent(cunits.Bits(1), 10000000, 0.0, 1, 0, "", 0), want: "3 months"},
		{input: makeTorrent(cunits.Bits(1), 100000000, 0.0, 1, 0, "", 0), want: "3 years"},
		{input: makeTorrent(cunits.Bits(1), -1, 0.0, 0, 0, "", 0), want: "Done"},
		{input: makeTorrent(cunits.Bits(1), 1, 0.0, 0, 0, "", 0), want: "Done"},
		{input: makeTorrent(cunits.Bits(1), -1, 0.0, 1, 0, "", 0), want: "Unknown"},
	}
	for _, tc := range tests {
		got := tc.input.Eta
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestNewFromTracker(t *testing.T) {
	tests := []test{
		{input: makeTorrent(cunits.Bits(1), 1, 0.0, 1, 0, "", 0), want: "UNK"},
		{input: makeTorrent(cunits.Bits(1), 1, 0.0, 1, 0, "http://foo-tracker", 0), want: "foo"},
	}
	for _, tc := range tests {
		got := tc.input.Trackershortname
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestNewFromProgress(t *testing.T) {
	tests := []test{
		{input: makeTorrent(cunits.Bits(8192), 1, 0.0, 256, 0, "", 0), want: int64(75)},
		{input: makeTorrent(cunits.Bits(8192), 1, 0.5, 256, 0, "", 0), want: int64(50)},
	}
	for _, tc := range tests {
		got := tc.input.Pct
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestNewFromRatio(t *testing.T) {
	tests := []test{
		// Nothing uploaded.
		{input: makeTorrent(cunits.Bits(8192), 1, 0.0, 256, 0, "", 0), want: 0.},
		// 512 out of 1024 uploaded.
		{input: makeTorrent(cunits.Bits(8192), 1, 0.5, 256, 512, "", 0), want: 0.5},
		// 512 uplaoded, 0 sized torrent (?!)
		{input: makeTorrent(cunits.Bits(0), 1, 0.5, 256, 512, "", 0), want: math.Inf(1)},
		// 0 uploaded, 0 sized torrent (?!)
		{input: makeTorrent(cunits.Bits(0), 1, 0.5, 256, 0, "", 0), want: math.NaN()},
	}
	for _, tc := range tests {
		got := tc.input.Ratio
		// Two NaNs aren't considered equal.
		if math.IsNaN(tc.want.(float64)) {
			if !math.IsNaN(got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		} else {
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		}
	}
}

func TestNewFromError(t *testing.T) {
	tests := []test{
		{input: makeTorrent(cunits.Bits(8192), 1, 0.0, 256, 0, "", 0), want: " "},
		{input: makeTorrent(cunits.Bits(8192), 1, 0.0, 256, 0, "", 1), want: "*"},
	}
	for _, tc := range tests {
		got := tc.input.Error
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
