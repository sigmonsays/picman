package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

var StateMask os.FileMode

func init() {
	StateMask = os.FileMode(0644)
}

func NewState() *State {
	ret := &State{}
	ret.Stat = &Stat{}
	ret.Date = &Date{}
	ret.Checksum = &Checksum{}
	ret.ExifData = &ExifData{}
	return ret
}

type State struct {

	// full path to the original filename we're importing
	OriginalFilename string

	// extension of filename (.jpg)
	Ext string

	// full path to the destination filenaem
	DestinationFilename string

	// the date the picture was taken
	Date *Date

	// data from stat
	Stat *Stat

	// checksum of the file
	Checksum *Checksum

	// exif data collected
	ExifData *ExifData

	// processing event logs
	Logs []string
}

type Stat struct {
	Size  int
	MTime time.Time
}

type Checksum struct {
	Sha256 string
}

type ExifData struct {
	Values map[string]string
}

type Date struct {
	Year, Month, Day     int
	Hour, Minute, Second int
}

func (me *State) Save(path string) error {
	// cs := sha256.New()
	// fmt.Fprintf(cs, me.OriginalFilename)
	// sha := cs.Sum(nil)
	// shaStr := hex.EncodeToString(sha)

	buf, err := json.MarshalIndent(me, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf, StateMask)
	if err != nil {
		return err
	}
	return nil
}
