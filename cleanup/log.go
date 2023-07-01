package cleanup

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("cleanup", func(newlog gologging.Logger) { log = newlog })
}
