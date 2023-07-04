package list

func (me *List) ProcessFile(srcdir string, statefile string, opts *Options) error {

	err := PrintState(srcdir, statefile, opts)
	if err != nil {
		log.Debugf("PrintState %s: %s", statefile, err)
	}

	return nil
}
