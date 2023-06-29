package core

type ImageProcessor interface {
	Run(*State) error
}
