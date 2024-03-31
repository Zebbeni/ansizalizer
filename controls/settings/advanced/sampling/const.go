package sampling

import "github.com/nfnt/resize"

var Functions = []resize.InterpolationFunction{
	resize.NearestNeighbor,
	resize.Bicubic,
	resize.Bilinear,
	resize.Lanczos2,
	resize.Lanczos3,
	resize.MitchellNetravali,
}

var nameMap = map[resize.InterpolationFunction]string{
	resize.NearestNeighbor:   "Nearest Neighbor",
	resize.Bicubic:           "Bicubic",
	resize.Bilinear:          "Bilinear",
	resize.Lanczos2:          "Lanczos2",
	resize.Lanczos3:          "Lanczos3",
	resize.MitchellNetravali: "MitchellNetravali",
}
