package autosort

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("autosort", func(newlog gologging.Logger) { log = newlog })
}
