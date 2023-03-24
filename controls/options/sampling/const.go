package sampling

import "github.com/nfnt/resize"

var Functions = []resize.InterpolationFunction{
	resize.NearestNeighbor,
	resize.Bilinear,
	resize.Bicubic,
	resize.MitchellNetravali,
	resize.Lanczos2,
	resize.Lanczos3,
}

var nameMap = map[resize.InterpolationFunction]string{
	resize.NearestNeighbor:   "Nearest Neighbor",
	resize.Bilinear:          "Bilinear",
	resize.Bicubic:           "Bicubic",
	resize.MitchellNetravali: "MitchellNetravali",
	resize.Lanczos2:          "Lanczos2",
	resize.Lanczos3:          "Lanczos3",
}
