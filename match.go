package main

import (
	"fmt"
	"image"

	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"

	"github.com/verocity-gaming/unitehud/dev"
	"github.com/verocity-gaming/unitehud/duplicate"
)

type match struct {
	image.Point
	template
}

func (m match) points(matrix2 gocv.Mat, img *image.RGBA) {
	results := make([]gocv.Mat, len(templates[points][m.Team.Name]))

	for i, pt := range templates[points][m.Team.Name] {
		mat := gocv.NewMat()
		defer mat.Close()

		results[i] = mat

		gocv.MatchTemplate(matrix2, pt.Mat, &mat, gocv.TmCcoeffNormed, mask)
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
	if value == 0 {
		log.Warn().Object("team", m.Team).Str("order", order).Msg("no value extracted")
	}

	region := m.Team.Region(matrix2)

	latest := duplicate.New(value, matrix2, region)

	dup := m.Team.Duplicate.Of(latest)
	if dup {
		log.Warn().Object("latest", latest).Object("last", m.Team.Duplicate).Msg("duplicate match")
	}

	m.Team.Duplicate.Close()
	m.Team.Duplicate = latest

	if record {
		dev.Capture(img, matrix2, m.Team.Name, order, dup, value)
		dev.Log(fmt.Sprintf("%s %d (duplicate: %t)", m.Team.Name, value, dup))
	}

	if !dup {
		go socket.Publish(m.Team, value)
	}
}

func (m match) process(matrix gocv.Mat, img *image.RGBA) {
	log.Info().Object("match", m).Int("cols", matrix.Cols()).Int("rows", matrix.Rows()).Msg("match found")

	switch m.category {
	case scored:
		m.points(matrix.Region(m.Team.Rectangle(m.Point)), img)
	case game:
		switch m.subcategory {
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
