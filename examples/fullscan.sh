#!/usr/bin/env sh

incoming=/data/Pictures

picman -l trace autosort \
    --source Pictures \
    --source-dir "$incoming" \
    --destination-dir /data/Pictures2 $@
