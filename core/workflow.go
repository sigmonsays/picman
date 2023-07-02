package core

import (
	"errors"
	"io/fs"
)

var (
	// stop processing current file and move on
	StopProcessing = errors.New("stop processing")

	// stop processing entire workflow and abort hard
	StopWorkflow = errors.New("stop workflow")

	// skip step and keep processing next step
	SkipStep = errors.New("skip workflow step")
)

type Workflow struct {

	// source flag given on cli (Phone10, Pixel7, etc)
	Source string

	// information about the image/video file
	Root     string      // --source-directory
	Fullpath string      // absolute path to image
	Info     fs.FileInfo // file stat data
	RelPath  string      // relative path to Root

	// destination path
	DestinationDir string // --destination-directory

	// instructs the CopyFile action to do nothing
	NoCopy bool
}
