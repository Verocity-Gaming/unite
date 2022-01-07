package duplicate

import (
	"image"
	"time"

	"gocv.io/x/gocv"

	"github.com/rs/zerolog"
	"github.com/verocity-gaming/unitehud/team"
)

type Duplicate struct {
	Value int
	team.Team
	time.Time
	gocv.Mat
}

func New(value int, team team.Team, mat gocv.Mat) *Duplicate {
	return &Duplicate{
		value,
		team,
		time.Now(),
		mat,
	}
}

func (d *Duplicate) Is(d2 *Duplicate) bool {
	if d == nil || d2 == nil {
		return false
	}

	if d.Empty() || d2.Empty() {
		return false
	}

	if d.Team.Name != d2.Team.Name {
		return false
	}

	if d.Value != d2.Value {
		return false
	}

	if d.Time.Sub(d2.Time) > time.Second*5 {
		return false
	}

	m := d2.Region()
	m2 := gocv.NewMat()

	gocv.MatchTemplate(d.Mat, m, &m2, gocv.TmCcoeffNormed, gocv.NewMat())
	_, maxc, _, _ := gocv.MinMaxLoc(m2)

	return maxc > 0.90
}

func (d *Duplicate) MarshalZerologObject(e *zerolog.Event) {
	e.Object("team", d.Team).Time("time", d.Time).Int("value", d.Value)
}

func (d *Duplicate) Region() gocv.Mat {
	if d.Team.Name == team.Self.Name {
		return d.Mat.Region(image.Rect(30, 20, 225, 60))
	}
	return d.Mat.Region(image.Rect(15, 30, 150, 60))
}
