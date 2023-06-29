package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var StateFileMask os.FileMode

func init() {
	StateFileMask = os.FileMode(0644)
}

type StateFile struct {

	// full path to the original filename we're importing
	OriginalFilename string

	// full path to the destination filenaem
	DestinationFilename string

	// the date the picture was taken
	Date *Date

	// checksum of the file
	Checksum *Checksum

	// exif data collected
	ExifData *ExifData

	// processing event logs
	Logs []string
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

func (me *StateFile) Save(path string) error {
	// cs := sha256.New()
	// fmt.Fprintf(cs, me.OriginalFilename)
	// sha := cs.Sum(nil)
	// shaStr := hex.EncodeToString(sha)

	buf, err := json.MarshalIndent(me, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf, StateFileMask)
	if err != nil {
		return err
	}
	return nil
}
