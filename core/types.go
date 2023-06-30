package core

import "os"

var DirMask = os.FileMode(0755)

type ImageProcessor interface {
	Run(*State) error
}
