package filter

import (
	"fmt"
	"net/url"
	"os"

	"github.com/shric/trpc/internal/fileutils"

	"github.com/shric/trpc/internal/torrent"

	"github.com/shric/trpc/internal/config"

	"github.com/shric/monkey/evaluator"
	"github.com/shric/monkey/lexer"
	"github.com/shric/monkey/object"
	"github.com/shric/monkey/parser"

	"github.com/hekmon/transmissionrpc"
)

// Options declares all the command line arguments for filtering torrents.
type Options struct {
	Filter      []string `short:"f" long:"filter" description:"apply filter expression" unquote:"false"`
	Incomplete  bool     `short:"i" long:"incomplete" description:"only incomplete torrents"`
	Active      bool     `short:"a" long:"active" description:"torrents currently uploading or downloading"`
	Name        string   `long:"name" description:"match a torrent name (regex)"`
	Tracker     string   `short:"t" long:"tracker" description:"match a tracker (regex)"`
	Error       string   `short:"e" long:"errors" description:"torrents with error matching string (regex)"`
	DownloadDir string   `long:"download-dir" description:"match on download directory"`
}

// Instance is used to hold all data required for a filter.
type Instance struct {
	conf        *config.Config
	expressions []string
	Args        []string
}

// New returns a new filter based on the options passed.
func New(opts Options, conf *config.Config) *Instance {
	expressions := opts.Filter
	args := make([]string, 0)

	if opts.Incomplete {
		expressions = append(expressions, "incomplete")
	}

	if opts.Active {
		expressions = append(expressions, "up > 0 || down > 0")
	}

	if opts.Tracker != "" {
		expressions = append(expressions, fmt.Sprintf("tracker ~ \"%s\"", opts.Tracker))
	}

	if opts.Error != "" {
		expressions = append(expressions, fmt.Sprintf("error ~ \"%s\"", opts.Error))
	}

	if opts.DownloadDir != "" {
		expressions = append(expressions, fmt.Sprintf("downloadDir == \"%s\"", fileutils.RealPath(opts.DownloadDir)))
	}

	if opts.Name != "" {
		expressions = append(expressions, fmt.Sprintf("name ~ \"%s\"", opts.Name))
	}

	filter := Instance{
		conf:        conf,
		expressions: expressions,
		Args:        args,
	}

	return &filter
}

func (f *Instance) envForTorrent(t *transmissionrpc.Torrent) *object.Environment {
	env := object.NewEnvironment()

	if *t.LeftUntilDone == 0 {
		env.Set("complete", evaluator.TRUE)
		env.Set("incomplete", evaluator.FALSE)
	} else {
		env.Set("complete", evaluator.FALSE)
		env.Set("incomplete", evaluator.TRUE)
	}

	env.Set("size", &object.Integer{Value: int64(t.SizeWhenDone.Byte())})

	trackers := make([]object.Object, len(t.Trackers))
	trackerStrings := make([]object.String, len(t.Trackers))

	for i, tracker := range t.Trackers {
		URL, err := url.Parse(tracker.Announce)
		if err != nil {
			continue
		}

		trackerStrings[i] = object.String{Value: URL.Hostname()}
		trackers[i] = &trackerStrings[i]
	}

	env.Set("trackers", &object.Array{Elements: trackers})
	env.Set("tracker", &object.String{Value: torrent.TrackerShortName(t, f.conf)})
	env.Set("down", &object.Integer{Value: *t.RateDownload})
	env.Set("up", &object.Integer{Value: *t.RateUpload})
	env.Set("age", &object.Integer{Value: torrent.Age(t)})
	env.Set("downloadDir", &object.String{Value: *t.DownloadDir})
	env.Set("priority", &object.String{Value: torrent.Priority(t)})
	env.Set("status", &object.String{Value: torrent.Status(t)})
	env.Set("name", &object.String{Value: *t.Name})

	if *t.Error != 0 {
		env.Set("error", &object.String{Value: *t.ErrorString})
	} else {
		env.Set("error", &object.String{Value: ""})
	}

	return env
}

// CheckFilter checks if the supplied torrent matches after filters.
func (f *Instance) CheckFilter(torrent *transmissionrpc.Torrent) bool {
	env := f.envForTorrent(torrent)

	for _, expr := range f.expressions {
		l := lexer.New(expr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			fmt.Fprintln(os.Stderr, "Filter parser error(s):")

			for _, msg := range p.Errors() {
				fmt.Fprintln(os.Stderr, "\t", msg)
			}

			fmt.Println(program.String())
			os.Exit(1)
		}

		result := evaluator.Eval(program, env)
		switch v := result.(type) {
		case *object.Boolean:
			if !v.Value {
				return false
			}

			continue
		case *object.Error:
			fmt.Fprintf(os.Stderr, "Invalid filter expression: %s\n", v.Message)
			os.Exit(1)
		default:
			fmt.Fprintf(os.Stderr, "Invalid filter expression: doesn't evaluate to boolean: %q\n", v)
			os.Exit(1)
		}
	}

	return true
}
