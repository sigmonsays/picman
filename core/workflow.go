package core

import (
	"errors"
	"io/fs"
)

var (
	// stop processing current file and move on
	StopProcessing = errors.New("stop processing")

	// stop processing entire workflow
	StopWorkflow = errors.New("stop workflow")

	// skip step and keep processing next step
	SkipStep = errors.New("skip workflow step")
)

type Workflow struct {

	// source flag given on cli
	Source string

	// information about the file
	Root     string
	Fullpath string
	Info     fs.FileInfo
	RelPath  string

	// destination path
	DestinationDir string
}
