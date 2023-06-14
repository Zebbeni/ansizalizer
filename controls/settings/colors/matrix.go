package colors

import (
	"github.com/makeworld-the-better-one/dither/v2"
)

type Matrix struct {
	Name   string
	Method dither.ErrorDiffusionMatrix
}

func getMatrixMenuItems() []Matrix {
	return []Matrix{
		Matrix{Name: "Simple2D", Method: dither.Simple2D},
		Matrix{Name: "FloydSteinberg", Method: dither.FloydSteinberg},
		Matrix{Name: "JarvisJudiceNinke", Method: dither.JarvisJudiceNinke},
		Matrix{Name: "Atkinson", Method: dither.Atkinson},
		Matrix{Name: "Stucki", Method: dither.Stucki},
		Matrix{Name: "Burkes", Method: dither.Burkes},
		Matrix{Name: "Sierra", Method: dither.Sierra},
		Matrix{Name: "StevenPigeon", Method: dither.StevenPigeon},
	}
}
