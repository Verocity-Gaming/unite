package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hybridgroup/mjpeg"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vova616/screenshot"
	"gocv.io/x/gocv"

	"github.com/verocity-gaming/unitehud/dev"
	"github.com/verocity-gaming/unitehud/duplicate"
	"github.com/verocity-gaming/unitehud/pipe"
	"github.com/verocity-gaming/unitehud/team"
)

// windows
// cls && go build && unitehud.exe -server -optimal

/*
	pieces={
		"0":{"file":"img/purple/points/point_1.png","point":"(45,10)","value":1},
		"1":{"file":"img/purple/points/point_0_alt.png","point":"(67,14)","value":0},
		"pieces":2
	}
		removed={
			"0":{"file":"img/purple/points/point_0_alt.png","point":"(67,14)","value":0},
			"pieces":1
	}
*/

type filter struct {
	team.Team
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

type match struct {
	image.Point
	template
}

var last *duplicate.Duplicate

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
	stream *mjpeg.Stream
	window *gocv.Window
	socket *pipe.Pipe

	mask    = gocv.NewMat()
	nomatch = match{}

	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}

	errWindowClosed = errors.New("window closed")
)

// Controllers.
var (
	sigq = make(chan os.Signal, 1)
)

// Consolidation.
var (
	skipped = 0
)

// Options.
var (
	cpu        = false
	record     = false
	delay      = time.Second / 2
	method     = gocv.TmCcoeffNormed
	acceptance = float32(0.91)
	//acceptance = float32(0.95)
	addr = ":17069"
)

func init() {
	flag.BoolVar(&record, "record", record, "record data such as images and logs for developer-specific debugging")
	flag.BoolVar(&cpu, "cpu", cpu, "decrease cpu load by optimizing image processing, with a potential impact on output")
	flag.DurationVar(&delay, "delay", delay, "interval to wait between capturing screen coordinates")
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

	for category := range filenames {
		for subcategory, filters := range filenames[category] {
			for _, filter := range filters {
				templates[category][filter.Team.Name] = append(templates[category][filter.Team.Name],
					template{
						filter,
						gocv.IMRead(filter.file, gocv.IMReadColor),
						1,
						category,
						subcategory,
					},
				)
			}
		}
	}

	for category := range templates {
		for _, templates := range templates[category] {
			for _, t := range templates {
				if t.Empty() {
					kill(fmt.Errorf("invalid scored template: %s (scale: %.2f)", t.file, t.scalar))
				}

				log.Debug().Object("template", t).Msg("score template loaded")
			}
		}
	}
}

func capture(category string) error {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return err
	}

	if cpu {
		// Cut the image in half.
		img.Rect.Max.Y /= 2
	}

	matrix, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return err
	}

	if matrix.Empty() {
		skipped++
		log.Warn().Int("skipped", skipped).Msg("skipped frame")
		return nil
	}

	m := matched(matrix, category)
	for _, match := range m {
		go match.process(matrix.Clone(), img)
	}

	return nil
}

func (match match) process(matrix gocv.Mat, img *image.RGBA) {
	log.Info().Object("match", match).Int("cols", matrix.Cols()).Int("rows", matrix.Rows()).Msg("match found")

	switch match.category {
	case scored:
		go match.points(matrix.Region(match.Team.Rectangle(match.Point)), img)
	case game:
		switch match.subcategory {
		case gameVS:
			socket.Clear()

			if record {
				dev.Start()
			}
		case gameEnd:
			if record {
				dev.End()
			}
		}
	}
}

func matched(matrix gocv.Mat, category string) []match {
	results := map[string][]gocv.Mat{}
	var matched []match

	for name, t := range templates[category] {
		results[name] = make([]gocv.Mat, len(templates[category][name]))

		for i, st := range t {
			mat := gocv.NewMat()
			defer mat.Close()

			gocv.MatchTemplate(matrix, st.Mat, &mat, method, mask)

			results[st.Team.Name][i] = mat
		}
	}

	// Iterate over resulting matches.
	for team, mats := range results {
		for i := range mats {
			if mats[i].Empty() {
				log.Warn().Str("filename", filenames[category][team][i].file).Msg("empty result")
				continue
			}

			_, maxc, _, maxp := gocv.MinMaxLoc(mats[i])
			if maxc >= acceptance {
				matched = append(matched, match{
					Point:    maxp,
					template: templates[category][team][i],
				})
			}
		}
	}

	return matched
}

func (m *match) points(matrix2 gocv.Mat, img *image.RGBA) {
	results := make([]gocv.Mat, len(templates[points][m.Team.Name]))

	for i, pt := range templates[points][m.Team.Name] {
		mat := gocv.NewMat()
		defer mat.Close()

		results[i] = mat

		gocv.MatchTemplate(matrix2, pt.Mat, &mat, method, mask)
	}

	pieces := pieces([]piece{})

	for i := range results {
		if results[i].Empty() {
			log.Warn().Str("filename", m.file).Msg("empty result")
			continue
		}

		_, maxc, _, maxp := gocv.MinMaxLoc(results[i])
		if maxc >= acceptance {
			pieces = append(pieces,
				piece{
					maxp,
					templates[points][m.Team.Name][i].filter,
				},
			)
		}
	}

	value, order := pieces.Int()

	latest := duplicate.New(value, m.Team, matrix2)

	if value == 0 {
		log.Warn().Object("team", m.Team).Str("order", order).Msg("no value extracted")
	}

	dup := last.Is(latest)
	if dup {
		log.Warn().Object("latest", latest).Object("last", last).Msg("duplicate match")
	}

	last = latest

	if record {
		dev.Capture(img, matrix2, latest.Team.Name, order, dup, latest.Value)
		dev.Log(fmt.Sprintf("%s %d (duplicate: %t)", latest.Name, latest.Value, dup))
	}

	if !dup {
		go socket.Publish(latest.Team, latest.Value)
	}
}

func kill(errs ...error) {
	for _, err := range errs {
		if err == errWindowClosed {
			continue
		}

		log.Error().Err(err).Send()
	}

	sigq <- os.Kill
}

func signals() {
	signal.Notify(sigq, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigq

	if s == os.Interrupt && window != nil {
		defer window.Close()
	}

	log.Info().Stringer("signal", s).Msg("closing...")
	os.Exit(1)
}

func main() {
	log.Info().
		Bool("record", record).
		Bool("cpu", cpu).
		Str("match", strconv.Itoa(int(acceptance*100))+"%").
		Str("addr", addr).Msg("unitehud")

	if record {
		err := dev.New()
		if err != nil {
			kill(err)
		}
	}

	socket = pipe.New(addr)

	for {
		for _, category := range []string{scored, game} {
			err := capture(category)
			if err != nil {
				kill(err)
			}
		}

		time.Sleep(delay)
	}
}

func (t template) scaled(scalar float64) template {
	img, err := t.ToImage()
	if err != nil {
		kill(err)
	}

	r := resize.Resize(uint(float64(t.Cols())*scalar), uint(float64(t.Rows())*scalar), img, resize.Lanczos3)

	m, err := gocv.ImageToMatRGB(r)
	if err != nil {
		kill(err)
	}

	return template{t.filter, m, scalar, t.category, t.subcategory}
}

func square(p image.Point, length, width int) image.Rectangle {
	return image.Rect(
		p.X,
		p.Y,
		p.X+length,
		p.Y+width,
	)
}
