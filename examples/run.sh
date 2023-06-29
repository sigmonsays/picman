#!/usr/bin/env sh


incoming=/data/Pictures-Android/AndroidDCIM/Camera

set -x

picman autosort \
    --source Phone10 \
    --source-dir "$incoming" \
    --destination-dir /data/Pictures
