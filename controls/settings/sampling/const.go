package sampling

import "github.com/nfnt/resize"

var Functions = []resize.InterpolationFunction{
	resize.Lanczos3,
	resize.Lanczos2,
	resize.NearestNeighbor,
	resize.Bilinear,
	resize.Bicubic,
	resize.MitchellNetravali,
}

var nameMap = map[resize.InterpolationFunction]string{
	resize.Lanczos3:          "Lanczos3",
	resize.Lanczos2:          "Lanczos2",
	resize.NearestNeighbor:   "Nearest Neighbor",
	resize.Bilinear:          "Bilinear",
	resize.Bicubic:           "Bicubic",
	resize.MitchellNetravali: "MitchellNetravali",
}
