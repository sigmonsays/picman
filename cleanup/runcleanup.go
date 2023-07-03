package cleanup

import (
	"os"
	"path/filepath"

	"github.com/sigmonsays/picman/core"
)

func NewReason(symbol string, msg string) *Reason {
	ret := &Reason{
		Symbol:  symbol,
		Message: msg,
	}
	return ret
}

type Reason struct {
	Symbol  string
	Message string
}

var (
	LoadError            = NewReason("LoadError", "load state error")
	DestinationEmpty     = NewReason("DestinationEmpty", "destination field is empty")
	NoDestination        = NewReason("NoDestination", "No destination file exists")
	DestinationIsNotFile = NewReason("DestinationIsNotFile", "destination is not a regular file")
	DoNotProcess         = NewReason("DoNotProcess", "do not process")
	FileTypeNotSupported = NewReason("FileTypeNotSupported", "file type not supported")
	MissingDate          = NewReason("MissingDate", "no date found")
	SourceMissing        = NewReason("SourceMissing", "source missing")
	//
	// LE DE ND DNF DNP FTNS MD SM
	// L D N F P NS D S
	//
	ReasonShortCodes = map[*Reason]string{
		LoadError:            "L",
		DestinationEmpty:     "E",
		NoDestination:        "N",
		DestinationIsNotFile: "F",
		DoNotProcess:         "P",
		FileTypeNotSupported: "T",
		MissingDate:          "D",
		SourceMissing:        "S",
	}
)

type Result struct {
	Reasons []*Reason
	Row     []string
}

func (me *Result) Finish() {
	expected := 5
	if len(me.Row) == expected {
		return
	}
	for i := 0; i < expected-len(me.Row); i++ {
		me.Row = append(me.Row, "\t")
	}
}
func (me *Result) WithReason(r *Reason) *Result {
	me.Reasons = append(me.Reasons, r)
	return me
}
func (me *Result) GetRow() []string {
	row := make([]string, 0)
	var reason *Reason
	if len(me.Reasons) > 0 {
		reason = me.Reasons[0]
	}
	if reason == nil {
		row = append(row, "OK")
	} else {
		codes := ""
		for _, r := range me.Reasons {
			shortcode, ok := ReasonShortCodes[r]
			if !ok {
				continue
			}
			codes += shortcode
		}
		row = append(row, codes)
	}
	return row
}

func RunCleanup(srcdir string, statefile string, opts *Options, stats *Stats) *Result {
	log.Tracef("")
	log.Tracef("start %s", statefile)
	stats.Processed++

	ret := &Result{}

	state := core.NewState()
	err := state.Load(statefile)
	if err != nil {
		return ret.WithReason(LoadError)
	}

	if state.DestinationFilename == "" {
		ret.WithReason(DestinationEmpty)
	}

	// make sure the source exists and is a regular file
	sourceExists := PathExistsRegularFile(state.OriginalFilename)
	if !sourceExists {
		ret.WithReason(SourceMissing)
	}

	// make sure the destination exists and is a regular file
	destExists := PathExistsRegularFile(state.OriginalFilename)
	if !destExists {
		ret.WithReason(NoDestination)
	}

	// check if our type is supported
	if core.IsFileExtSupported(state.Ext) != nil {
		ret.WithReason(FileTypeNotSupported)
	}

	// make sure we have a date
	if state.Date == nil || (state.Date.Year == 0 || state.Date.Month == 0) {
		ret.WithReason(MissingDate)
	}

	// output the row
	// - reason
	// - ext
	// - srcpath (relative)
	// - dstpath (absolute)
	row := ret.GetRow()

	if state.Ext == "" {
		row = append(row, "-")
	} else {
		row = append(row, state.Ext)

	}

	// src path
	srcrel, _ := filepath.Rel(srcdir, state.OriginalFilename)
	row = append(row, srcrel)

	// dest path
	row = append(row, state.DestinationFilename)
	ret.Row = row

	// todo: if DoNotProcess is set we just delete the metadata file
	// maybe we need a --delete flag

	return ret
}

// return true if path exists and is a regular file
func PathExistsRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.Mode().IsRegular() == false {
		return false
	}
	return true
}
