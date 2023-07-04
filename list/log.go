package list

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("list", func(newlog gologging.Logger) { log = newlog })
}
