package core

import (
	"encoding/json"
	"fmt"
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

	// source device of image
	// - Phone10 for phone
	Source string

	// do not process
	DoNotProcess bool

	// true if file copied to final destination filename
	FileCopied bool

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

type Date struct {
	Year, Month, Day     int
	Hour, Minute, Second int
}

func (me *State) StopProcessing(s string, args ...interface{}) error {
	me.DoNotProcess = true
	me.Logf(s, args...)
	return StopProcessing
}

func (me *State) Load(path string) error {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, me)
	if err != nil {
		return err
	}
	return nil

}
func (me *State) Save(path string) error {

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

func (me *State) Logf(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	me.Logs = append(me.Logs, msg)
}
