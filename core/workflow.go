package core

import (
	"io/fs"
)

type Workflow struct {
	// information about the file
	Root     string
	Fullpath string
	Info     fs.FileInfo
	RelPath  string

	// destination path
	DestinationDir string
}
