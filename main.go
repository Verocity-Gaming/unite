package main

import (
	"flag"
	"image"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vova616/screenshot"
	"gocv.io/x/gocv"

	"github.com/verocity-gaming/unitehud/dev"
	"github.com/verocity-gaming/unitehud/pipe"
	"github.com/verocity-gaming/unitehud/team"
)

// windows
// cls && go build && unitehud.exe -server -optimal

type filter struct {
	*team.Team
	file  string
	value int
}

type template struct {
	filter
	gocv.Mat
	scalar      float64
	category    string
	subcategory string
}

// Categories
const (
	game    = "game"
	gameVS  = "vs"
	gameEnd = "end"
	scored  = "scored"
	points  = "points"
)

// Coadjutants.
var (
	socket *pipe.Pipe
	mask   = gocv.NewMat()
	rect   = image.Rect(640, 0, 1280, 500)
	sigq   = make(chan os.Signal, 1)
)

// Options.
var (
	record     = false
	acceptance = float32(0.91)
	//acceptance = float32(0.95)
	addr = ":17069"
)

var workers = map[string]int{
	team.None.Name:   1,
	team.Self.Name:   4,
	team.Purple.Name: 1,
	team.Orange.Name: 1,
}

var distribq = map[string]chan *image.RGBA{
	team.None.Name:   make(chan *image.RGBA),
	team.Self.Name:   make(chan *image.RGBA),
	team.Purple.Name: make(chan *image.RGBA),
	team.Orange.Name: make(chan *image.RGBA),
}

var imgq = map[string]chan *image.RGBA{
	team.None.Name:   make(chan *image.RGBA),
	team.Self.Name:   make(chan *image.RGBA),
	team.Purple.Name: make(chan *image.RGBA),
	team.Orange.Name: make(chan *image.RGBA),
}

func init() {
	flag.BoolVar(&record, "record", record, "record data such as images and logs for developer-specific debugging")
	flag.StringVar(&addr, "addr", addr, "http/websocket serve address")
	avg := flag.Float64("match", float64(acceptance)*100, `0-100% certainty when processing score values`)
	level := flag.String("v", zerolog.LevelDebugValue, "log level (panic, fatal, error, warn, info, debug)")
	flag.Parse()

	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.Stamp,
		},
	).With().Timestamp().Logger()

	acceptance = float32(*avg) / 100

	go signals()

	lvl, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Logger = log.Logger.Level(lvl)

	load()
}

func capture(name string) {
	for {
		img, err := screenshot.CaptureRect(rect)
		if err != nil {
			kill(err)
		}

		select {
		case imgq[name] <- img:
		default:
		}

		time.Sleep(team.Delay(name))
	}
}

func loop(t []template, imgq chan *image.RGBA) {
	runtime.LockOSThread()

	for img := range imgq {
		matrix, err := gocv.ImageToMatRGB(img)
		if err != nil {
			kill(err)
		}

		matches(matrix, img, t)

		matrix.Close()
	}
}

func matches(matrix gocv.Mat, img *image.RGBA, t []template) {
	results := make([]gocv.Mat, len(t))

	for i, template := range t {
		results[i] = gocv.NewMat()
		defer results[i].Close()

		gocv.MatchTemplate(matrix, template.Mat, &results[i], gocv.TmCcoeffNormed, mask)
	}

	for i, mat := range results {
		if mat.Empty() {
			log.Warn().Str("filename", t[i].file).Msg("empty result")
			continue
		}

		_, maxc, _, maxp := gocv.MinMaxLoc(mat)
		if maxc >= acceptance {
			match{
				Point:    maxp,
				template: t[i],
			}.process(matrix, img)
		}
	}
}

func kill(errs ...error) {
	for _, err := range errs {
		log.Error().Err(err).Send()
		time.Sleep(time.Millisecond)
	}

	sigq <- os.Kill
}

func signals() {
	signal.Notify(sigq, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigq

	log.Info().Stringer("signal", s).Msg("closing...")

	os.Exit(1)
}

func main() {
	log.Info().
		Bool("record", record).
		Str("match", strconv.Itoa(int(acceptance*100))+"%").
		Str("addr", addr).Msg("unitehud")

	if record {
		err := dev.New()
		if err != nil {
			kill(err)
		}
	}

	socket = pipe.New(addr)

	for category := range templates {
		if category == points {
			continue
		}

		for name := range templates[category] {
			for i := 0; i < workers[name]; i++ {
				go capture(name)
				go loop(templates[category][name], imgq[name])
			}
		}
	}

	signals()
}
