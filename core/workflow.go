package core

import (
	"io/fs"
)

type Workflow struct {
	Root     string
	Fullpath string
	Info     fs.FileInfo
	RelPath  string
}
