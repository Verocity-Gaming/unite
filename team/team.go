package team

import (
	"image"
	"image/color"

	"github.com/rs/zerolog"
)

var (
	black  = color.RGBA{0, 0, 0, 255}
	orange = color.RGBA{255, 165, 0, 255}
	purple = color.RGBA{83, 94, 255, 255}
	// white  = color.RGBA{255, 255, 255, 255}
)

// Team represents a team side in Pokemon Unite.
type Team struct {
	Name       string `json:"name"`
	color.RGBA `json:"-"`
}

var (
	// Orange represents the standard Team for the Orange side.
	Orange = Team{
		Name: "orange",
		RGBA: orange,
	}

	// Purple represents the standard Team for the Purple side.
	Purple = Team{
		Name: "purple",
		RGBA: purple,
	}

	// Self represents a wrapper Team for the Purple side.
	Self = Team{
		Name: "self",
		RGBA: purple,
	}

	// None represents the default case for an unidentifiable side.
	None = Team{
		Name: "none",
		RGBA: black,
	}
)

func (t Team) Rectangle(p image.Point) image.Rectangle {
	if t.Name == Self.Name {
		return image.Rect(p.X-250, p.Y-50, p.X+300, p.Y+100)
	}
	return image.Rect(p.X-100, p.Y-30, p.X+150, p.Y+75)
}

func (t Team) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", t.Name)
}
