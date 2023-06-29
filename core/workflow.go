package core

import (
	"io/fs"
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
